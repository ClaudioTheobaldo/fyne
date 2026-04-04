package gl

import (
	"image"
	"image/color"
	"math"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	paint "fyne.io/fyne/v2/internal/painter"
)

// BenchImageDrawNs accumulates nanoseconds spent in drawImage (for benchmarking).
var BenchImageDrawNs atomic.Int64

// BenchStreamDrawNs accumulates nanoseconds spent in drawStreamingImage (for benchmarking).
var BenchStreamDrawNs atomic.Int64

// BenchImageDrawCount counts drawImage invocations (for benchmarking).
var BenchImageDrawCount atomic.Int64

// BenchStreamDrawCount counts drawStreamingImage invocations (for benchmarking).
var BenchStreamDrawCount atomic.Int64

// edgeSoftness returns the SDF edge softness value based on the current AA mode.
// When AA is Off, returns 0 (hard edges); otherwise returns 1.0.
func (p *painter) edgeSoftness() float32 {
	if !p.aaMode.NeedsSDF() {
		return 0.0
	}
	return 1.0
}

func (p *painter) createBuffer(size int) Buffer {
	vbo := p.ctx.CreateBuffer()
	p.logError()
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	p.ctx.BufferData(arrayBuffer, make([]float32, size), staticDraw)
	p.logError()
	return vbo
}

func (p *painter) drawCircle(circle *canvas.Circle, pos fyne.Position, frame fyne.Size) {
	radius := paint.GetMaximumRadius(circle.Size())
	program := p.roundRectangleProgram

	// Vertex: BEG
	bounds, points := p.vecSquareCoords(pos, circle, frame)
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(circle.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width_half", strokeWidthScaled*0.5)

	rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
	rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
	p.SetUniform2f(program, "rect_size_half", rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

	radiusScaled := roundToPixel(radius*p.pixScale, 1.0)
	p.SetUniform4f(program, "radius", radiusScaled, radiusScaled, radiusScaled, radiusScaled)

	r, g, b, a := getFragmentColor(circle.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := circle.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawGradient(o fyne.CanvasObject, texCreator func(fyne.CanvasObject) Texture, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(o, texCreator, pos, o.Size(), frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *painter) drawImage(img *canvas.Image, pos fyne.Position, frame fyne.Size) {
	t0 := time.Now()
	p.drawTextureWithDetails(img, p.newGlImageTexture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0)
	BenchImageDrawNs.Add(time.Since(t0).Nanoseconds())
	BenchImageDrawCount.Add(1)
}

func (p *painter) drawLine(line *canvas.Line, pos fyne.Position, frame fyne.Size) {
	if line.StrokeColor == color.Transparent || line.StrokeColor == nil || line.StrokeWidth == 0 {
		return
	}
	featherVal := float32(0.5)
	if !p.aaMode.NeedsSDF() {
		featherVal = 0.0
	}
	points, halfWidth, feather := p.lineCoords(pos, line.Position1, line.Position2, line.StrokeWidth, featherVal, frame)
	p.ctx.UseProgram(p.lineProgram.ref)
	p.updateBuffer(p.lineProgram.buff, points)
	p.UpdateVertexArray(p.lineProgram, "vert", 2, 4, 0)
	p.UpdateVertexArray(p.lineProgram, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	r, g, b, a := getFragmentColor(line.StrokeColor)
	p.SetUniform4f(p.lineProgram, "color", r, g, b, a)

	p.SetUniform1f(p.lineProgram, "lineWidth", halfWidth)

	p.SetUniform1f(p.lineProgram, "feather", feather)

	p.ctx.DrawArrays(triangles, 0, 6)
	p.logError()
}

func (p *painter) drawObject(o fyne.CanvasObject, pos fyne.Position, frame fyne.Size) {
	if holder, ok := o.(effectHolder); ok && holder.HasEffects() {
		p.drawObjectWithEffects(o, pos, frame)
		return
	}
	switch obj := o.(type) {
	case *canvas.StreamingImage:
		p.drawStreamingImage(obj, pos, frame)
	case *canvas.Circle:
		p.drawCircle(obj, pos, frame)
	case *canvas.Line:
		p.drawLine(obj, pos, frame)
	case *canvas.Image:
		p.drawImage(obj, pos, frame)
	case *canvas.Raster:
		p.drawRaster(obj, pos, frame)
	case *canvas.Rectangle:
		p.drawRectangle(obj, pos, frame)
	case *canvas.Text:
		p.drawText(obj, pos, frame)
	case *canvas.LinearGradient:
		p.drawGradient(obj, p.newGlLinearGradientTexture, pos, frame)
	case *canvas.RadialGradient:
		p.drawGradient(obj, p.newGlRadialGradientTexture, pos, frame)
	case *canvas.Polygon:
		p.drawPolygon(obj, pos, frame)
	case *canvas.Arc:
		p.drawArc(obj, pos, frame)
	case *canvas.ShaderRect:
		p.drawShaderRect(obj, pos, frame)
	}
}

func (p *painter) drawRaster(img *canvas.Raster, pos fyne.Position, frame fyne.Size) {
	p.drawTextureWithDetails(img, p.newGlRasterTexture, pos, img.Size(), frame, canvas.ImageFillStretch, float32(img.Alpha()), 0)
}

func (p *painter) drawRectangle(rect *canvas.Rectangle, pos fyne.Position, frame fyne.Size) {
	topRightRadius := paint.GetCornerRadius(rect.TopRightCornerRadius, rect.CornerRadius)
	topLeftRadius := paint.GetCornerRadius(rect.TopLeftCornerRadius, rect.CornerRadius)
	bottomRightRadius := paint.GetCornerRadius(rect.BottomRightCornerRadius, rect.CornerRadius)
	bottomLeftRadius := paint.GetCornerRadius(rect.BottomLeftCornerRadius, rect.CornerRadius)
	p.drawOblong(rect, rect.FillColor, rect.StrokeColor, rect.StrokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, rect.Aspect, pos, frame)
}

func (p *painter) drawOblong(obj fyne.CanvasObject, fill, stroke color.Color, strokeWidth, topRightRadius, topLeftRadius, bottomRightRadius, bottomLeftRadius, aspect float32, pos fyne.Position, frame fyne.Size) {
	if (fill == color.Transparent || fill == nil) && (stroke == color.Transparent || stroke == nil || strokeWidth == 0) {
		return
	}

	roundedCorners := topRightRadius != 0 || topLeftRadius != 0 || bottomRightRadius != 0 || bottomLeftRadius != 0
	var program ProgramState
	if roundedCorners {
		program = p.roundRectangleProgram
	} else {
		program = p.rectangleProgram
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, obj, frame, aspect)
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	strokeWidthScaled := roundToPixel(strokeWidth*p.pixScale, 1.0)
	if roundedCorners {
		p.SetUniform1f(program, "stroke_width_half", strokeWidthScaled*0.5)

		rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
		rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
		p.SetUniform2f(program, "rect_size_half", rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

		// the maximum possible corner radii for a circular shape, calculated taking into account the rect coords with aspect ratio
		size := fyne.NewSize(bounds[2]-bounds[0], bounds[3]-bounds[1])
		topRightRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(topRightRadius, topLeftRadius, bottomRightRadius, size)*p.pixScale,
			1.0,
		)
		topLeftRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(topLeftRadius, topRightRadius, bottomLeftRadius, size)*p.pixScale,
			1.0,
		)
		bottomRightRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(bottomRightRadius, bottomLeftRadius, topRightRadius, size)*p.pixScale,
			1.0,
		)
		bottomLeftRadiusScaled := roundToPixel(
			paint.GetMaximumCornerRadius(bottomLeftRadius, bottomRightRadius, topLeftRadius, size)*p.pixScale,
			1.0,
		)
		p.SetUniform4f(program, "radius", topRightRadiusScaled, bottomRightRadiusScaled, topLeftRadiusScaled, bottomLeftRadiusScaled)

		edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
		p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)
	} else {
		p.SetUniform1f(program, "stroke_width", strokeWidthScaled)
	}

	r, g, b, a := getFragmentColor(fill)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := stroke
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)
	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawShaderRect(rect *canvas.ShaderRect, pos fyne.Position, frame fyne.Size) {
	if rect.FragmentShader == "" {
		return
	}
	fill := rect.FillColor
	stroke := rect.StrokeColor
	if (fill == color.Transparent || fill == nil) && (stroke == color.Transparent || stroke == nil || rect.StrokeWidth == 0) {
		return
	}

	// Select shader source based on platform
	fragSrc := rect.FragmentShader
	if isGLES() && rect.FragmentShaderES != "" {
		fragSrc = rect.FragmentShaderES
	}

	ps := p.getOrCompileShaderProgram(fragSrc)
	if ps == nil {
		return // compilation failed
	}

	// Discover custom uniforms on first use
	if len(rect.Uniforms) > 0 {
		p.discoverCustomUniforms(ps, ps.ref, rect.Uniforms)
	}

	// Vertex setup (same as drawOblong)
	bounds, points := p.vecRectCoords(pos, rect, frame, 0)
	p.ctx.UseProgram(ps.ref)
	p.updateBuffer(ps.buff, points)
	p.UpdateVertexArray(*ps, "vert", 2, 4, 0)
	p.UpdateVertexArray(*ps, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()

	// Standard uniforms
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(*ps, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(*ps, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	topRightRadius := paint.GetCornerRadius(rect.TopRightCornerRadius, rect.CornerRadius)
	topLeftRadius := paint.GetCornerRadius(rect.TopLeftCornerRadius, rect.CornerRadius)
	bottomRightRadius := paint.GetCornerRadius(rect.BottomRightCornerRadius, rect.CornerRadius)
	bottomLeftRadius := paint.GetCornerRadius(rect.BottomLeftCornerRadius, rect.CornerRadius)

	strokeWidthScaled := roundToPixel(rect.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(*ps, "stroke_width_half", strokeWidthScaled*0.5)

	rectSizeWidthScaled := x2Scaled - x1Scaled - strokeWidthScaled
	rectSizeHeightScaled := y2Scaled - y1Scaled - strokeWidthScaled
	p.SetUniform2f(*ps, "rect_size_half", rectSizeWidthScaled*0.5, rectSizeHeightScaled*0.5)

	size := fyne.NewSize(bounds[2]-bounds[0], bounds[3]-bounds[1])
	topRightRadiusScaled := roundToPixel(paint.GetMaximumCornerRadius(topRightRadius, topLeftRadius, bottomRightRadius, size)*p.pixScale, 1.0)
	topLeftRadiusScaled := roundToPixel(paint.GetMaximumCornerRadius(topLeftRadius, topRightRadius, bottomLeftRadius, size)*p.pixScale, 1.0)
	bottomRightRadiusScaled := roundToPixel(paint.GetMaximumCornerRadius(bottomRightRadius, bottomLeftRadius, topRightRadius, size)*p.pixScale, 1.0)
	bottomLeftRadiusScaled := roundToPixel(paint.GetMaximumCornerRadius(bottomLeftRadius, bottomRightRadius, topLeftRadius, size)*p.pixScale, 1.0)
	p.SetUniform4f(*ps, "radius", topRightRadiusScaled, bottomRightRadiusScaled, topLeftRadiusScaled, bottomLeftRadiusScaled)

	edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
	p.SetUniform1f(*ps, "edge_softness", edgeSoftnessScaled)

	r, g, b, a := getFragmentColor(fill)
	p.SetUniform4f(*ps, "fill_color", r, g, b, a)

	strokeColor := stroke
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(*ps, "stroke_color", r, g, b, a)

	// User custom uniforms
	for name, vals := range rect.Uniforms {
		switch len(vals) {
		case 1:
			p.SetUniform1f(*ps, name, vals[0])
		case 2:
			p.SetUniform2f(*ps, name, vals[0], vals[1])
		case 4:
			p.SetUniform4f(*ps, name, vals[0], vals[1], vals[2], vals[3])
		}
	}
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawPolygon(polygon *canvas.Polygon, pos fyne.Position, frame fyne.Size) {
	if ((polygon.FillColor == color.Transparent || polygon.FillColor == nil) && (polygon.StrokeColor == color.Transparent || polygon.StrokeColor == nil || polygon.StrokeWidth == 0)) || polygon.Sides < 3 {
		return
	}
	size := polygon.Size()

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, polygon, frame, 0.0)
	program := p.polygonProgram
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)

	outerRadius := fyne.Min(size.Width, size.Height) / 2
	outerRadiusScaled := roundToPixel(outerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "outer_radius", outerRadiusScaled)

	p.SetUniform1f(program, "angle", polygon.Angle)
	p.SetUniform1f(program, "sides", float32(polygon.Sides))

	cornerRadius := fyne.Min(paint.GetMaximumRadius(size), polygon.CornerRadius)
	cornerRadiusScaled := roundToPixel(cornerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "corner_radius", cornerRadiusScaled)

	strokeWidthScaled := roundToPixel(polygon.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width", strokeWidthScaled)

	r, g, b, a := getFragmentColor(polygon.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := polygon.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawArc(arc *canvas.Arc, pos fyne.Position, frame fyne.Size) {
	if ((arc.FillColor == color.Transparent || arc.FillColor == nil) && (arc.StrokeColor == color.Transparent || arc.StrokeColor == nil || arc.StrokeWidth == 0)) || arc.StartAngle == arc.EndAngle {
		return
	}

	// Vertex: BEG
	bounds, points := p.vecRectCoords(pos, arc, frame, 0.0)
	program := p.arcProgram
	p.ctx.UseProgram(program.ref)
	p.updateBuffer(program.buff, points)
	p.UpdateVertexArray(program, "vert", 2, 4, 0)
	p.UpdateVertexArray(program, "normal", 2, 4, 2)

	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.logError()
	// Vertex: END

	// Fragment: BEG
	frameWidthScaled, frameHeightScaled := p.scaleFrameSize(frame)
	p.SetUniform2f(program, "frame_size", frameWidthScaled, frameHeightScaled)

	x1Scaled, x2Scaled, y1Scaled, y2Scaled := p.scaleRectCoords(bounds[0], bounds[2], bounds[1], bounds[3])
	p.SetUniform4f(program, "rect_coords", x1Scaled, x2Scaled, y1Scaled, y2Scaled)

	edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
	p.SetUniform1f(program, "edge_softness", edgeSoftnessScaled)

	outerRadius := fyne.Min(arc.Size().Width, arc.Size().Height) / 2
	outerRadiusScaled := roundToPixel(outerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "outer_radius", outerRadiusScaled)

	innerRadius := outerRadius * float32(math.Min(1.0, math.Max(0.0, float64(arc.CutoutRatio))))
	innerRadiusScaled := roundToPixel(innerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "inner_radius", innerRadiusScaled)

	startAngle, endAngle := paint.NormalizeArcAngles(arc.StartAngle, arc.EndAngle)
	p.SetUniform1f(program, "start_angle", startAngle)
	p.SetUniform1f(program, "end_angle", endAngle)

	cornerRadius := fyne.Min(paint.GetMaximumRadiusArc(outerRadius, innerRadius, arc.EndAngle-arc.StartAngle), arc.CornerRadius)
	cornerRadiusScaled := roundToPixel(cornerRadius*p.pixScale, 1.0)
	p.SetUniform1f(program, "corner_radius", cornerRadiusScaled)

	strokeWidthScaled := roundToPixel(arc.StrokeWidth*p.pixScale, 1.0)
	p.SetUniform1f(program, "stroke_width", strokeWidthScaled)

	r, g, b, a := getFragmentColor(arc.FillColor)
	p.SetUniform4f(program, "fill_color", r, g, b, a)

	strokeColor := arc.StrokeColor
	if strokeColor == nil {
		strokeColor = color.Transparent
	}
	r, g, b, a = getFragmentColor(strokeColor)
	p.SetUniform4f(program, "stroke_color", r, g, b, a)

	p.logError()
	// Fragment: END

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawText(text *canvas.Text, pos fyne.Position, frame fyne.Size) {
	if text.Text == "" || text.Text == " " {
		return
	}

	size := text.MinSize()
	containerSize := text.Size()
	switch text.Alignment {
	case fyne.TextAlignTrailing:
		pos = fyne.NewPos(pos.X+containerSize.Width-size.Width, pos.Y)
	case fyne.TextAlignCenter:
		pos = fyne.NewPos(pos.X+(containerSize.Width-size.Width)/2, pos.Y)
	}

	if containerSize.Height > size.Height {
		pos = fyne.NewPos(pos.X, pos.Y+(containerSize.Height-size.Height)/2)
	}

	// text size is sensitive to position on screen
	size.Width = roundToPixel(size.Width, p.pixScale)
	size.Height = roundToPixel(size.Height, p.pixScale)
	size.Width += roundToPixel(paint.VectorPad(text), p.pixScale)
	p.drawTextureWithDetails(text, p.newGlTextTexture, pos, size, frame, canvas.ImageFillStretch, 1.0, 0)
}

func (p *painter) drawStreamingImage(img *canvas.StreamingImage, pos fyne.Position, frame fyne.Size) {
	// Check for new raw frame OR new RGBA frame.
	// If an RGBA frame (UpdateFrame) is pending, it takes priority over stale raw PBO state.
	// This handles the YUV→RGBA transition (e.g. gap mode showing a blue placeholder).
	rawFrame := img.ConsumePendingRawFrame()
	newRGBAFrame := img.ConsumePendingFrame()

	// Determine which rendering path to use:
	// 1. New RGBA frame pending (UpdateFrame) → RGBA path, clear stale PBO
	// 2. New raw frame pending (UpdateRawFrame) or existing PBO state → raw path
	// 3. Neither pending → use whichever cached state exists (raw PBO or RGBA texture)
	hasRawState := p.rawPBOStates != nil && p.rawPBOStates[img] != nil && p.rawPBOStates[img].ready

	if newRGBAFrame != nil && rawFrame == nil {
		// RGBA frame takes priority — clear stale raw PBO state so subsequent
		// repaints without new data use the RGBA texture cache, not stale PBO.
		if p.rawPBOStates != nil {
			delete(p.rawPBOStates, img)
		}
		hasRawState = false
	}

	if rawFrame != nil || hasRawState {
		p.drawStreamingImageRaw(img, rawFrame, pos, frame)
		return
	}

	// RGBA path.
	t0 := time.Now()
	defer func() {
		BenchStreamDrawNs.Add(time.Since(t0).Nanoseconds())
		BenchStreamDrawCount.Add(1)
	}()

	var texture Texture
	if newRGBAFrame != nil {
		texture = p.uploadStreamingFrame(img, newRGBAFrame)
	} else if existingTex, cached := cache.GetTexture(img); cached {
		texture = Texture(existingTex)
	} else {
		// Nothing to draw — no pending frame and no cached texture.
		// This only happens before the very first frame arrives.
		// Do NOT clear the screen — just skip drawing this object entirely
		// so the previous content (if any) remains visible via the framebuffer.
		return
	}

	p.drawQuadWithTexture(texture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0, 0)
}

// drawStreamingImageRaw draws a StreamingImage using the generic RawFrame path.
// rawFrame may be nil when repainting with no new data (uses cached GPU state).
func (p *painter) drawStreamingImageRaw(img *canvas.StreamingImage, rawFrame *canvas.RawFrame, pos fyne.Position, frame fyne.Size) {
	t0 := time.Now()
	defer func() {
		BenchStreamDrawNs.Add(time.Since(t0).Nanoseconds())
		BenchStreamDrawCount.Add(1)
	}()

	desc := formatDescriptorFor(img.PixelFormat)
	colorRange := img.ColorRange
	if colorRange == 0 && desc.defaultRange != 0 {
		colorRange = desc.defaultRange
	}

	switch desc.category {
	case categorySimple:
		var texture Texture
		if rawFrame != nil {
			stride := rawFrame.Strides[0]
			if stride == 0 {
				stride = rawFrame.Width * 4
			}
			rgba := &image.RGBA{Pix: rawFrame.Data[0], Stride: stride,
				Rect: image.Rect(0, 0, rawFrame.Width, rawFrame.Height)}
			texture = p.uploadStreamingFrame(img, rgba)
		} else if existingTex, cached := cache.GetTexture(img); cached {
			texture = Texture(existingTex)
		} else {
			return
		}
		p.drawQuadWithTexture(texture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0, 0)

	case categoryCPUConvert:
		var texture Texture
		if rawFrame != nil {
			rgba := convertRawToRGBA(img.PixelFormat, rawFrame)
			if rgba == nil {
				return
			}
			texture = p.uploadStreamingFrame(img, rgba)
		} else if existingTex, cached := cache.GetTexture(img); cached {
			texture = Texture(existingTex)
		} else {
			return
		}
		p.drawQuadWithTexture(texture, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()), 0, 0)

	case categoryYUVPlanar:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, desc.bitDepth)
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithYUVPlanar(textures, cm, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryYUVAPlanar:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, desc.bitDepth)
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithYUVAPlanar(textures, cm, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryNVSemiplanar:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, 8)
		swapUV := float32(0)
		if desc.swapUV {
			swapUV = 1
		}
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithNVSemiplanar(textures, swapUV, cm, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryPackedYUV422:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, 8)
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		texW, _ := img.TextureSize()
		texPackedWidth := float32(texW) // actual width set after upload
		if rawFrame != nil {
			texPackedWidth = float32(rawFrame.Width / desc.chromaWDiv[0])
		}
		p.drawQuadWithPackedYUV422(textures[0], float32(desc.packingMode), texPackedWidth, cm,
			pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryGrayscale:
		bitMax := float32(1.0)
		if desc.bitDepth == 16 {
			bitMax = 65535.0
		}
		hasAlpha := float32(0)
		if desc.hasAlpha {
			hasAlpha = 1
		}
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithGrayscale(textures[0], bitMax, hasAlpha, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryHiBitPlanar:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, desc.bitDepth)
		maxVal := (1 << uint(desc.bitDepth)) - 1
		bitMax := float32(maxVal)
		hasAlpha := float32(0)
		if desc.hasAlpha {
			hasAlpha = 1
		}
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithHiBitPlanar(textures, bitMax, hasAlpha, cm, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))

	case categoryHiBitNV:
		cm := ComputeColorMatrix(img.ColorSpace, colorRange, desc.bitDepth)
		maxValNV := (1 << uint(desc.bitDepth)) - 1
		bitMax := float32(maxValNV)
		swapUV := float32(0)
		if desc.swapUV {
			swapUV = 1
		}
		var textures [4]Texture
		if rawFrame != nil {
			textures = p.uploadStreamingRawFrame(img, rawFrame, desc)
		} else if state := p.rawPBOStates[img]; state != nil {
			textures = state.textures
		} else {
			return
		}
		p.drawQuadWithHiBitNV(textures, bitMax, swapUV, cm, pos, img.Size(), frame, img.FillMode, float32(img.Alpha()))
	}
}

// setColorMatrixUniforms uploads the 4 color-matrix vec3 uniforms.
func (p *painter) setColorMatrixUniforms(prog ProgramState, cm ColorMatrix) {
	p.SetUniform3f(prog, "colorRow0", cm.Row0[0], cm.Row0[1], cm.Row0[2])
	p.SetUniform3f(prog, "colorRow1", cm.Row1[0], cm.Row1[1], cm.Row1[2])
	p.SetUniform3f(prog, "colorRow2", cm.Row2[0], cm.Row2[1], cm.Row2[2])
	p.SetUniform3f(prog, "colorOffset", cm.Offset[0], cm.Offset[1], cm.Offset[2])
}

// setStreamProgramBase sets the geometry/alpha uniforms shared by all streaming programs.
func (p *painter) setStreamProgramBase(prog ProgramState, points []float32, inner fyne.Size, insets [4]float32, alpha float32) {
	p.updateBuffer(prog.buff, points)
	p.UpdateVertexArray(prog, "vert", 3, 5, 0)
	p.UpdateVertexArray(prog, "vertTexCoord", 2, 5, 3)
	p.SetUniform1f(prog, "cornerRadius", 0)
	p.SetUniform2f(prog, "size", inner.Width*p.pixScale, inner.Height*p.pixScale)
	p.SetUniform4f(prog, "inset", insets[0], insets[1], insets[2], insets[3])
	p.SetUniform1f(prog, "alpha", alpha)
}

func (p *painter) drawQuadWithYUVPlanar(textures [4]Texture, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.yuvPlanarProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, textures[0])
	p.ctx.ActiveTexture(texture1)
	p.ctx.BindTexture(texture2D, textures[1])
	p.ctx.ActiveTexture(texture2)
	p.ctx.BindTexture(texture2D, textures[2])
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithYUVAPlanar(textures [4]Texture, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.yuvaPlanarProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, textures[0])
	p.ctx.ActiveTexture(texture1)
	p.ctx.BindTexture(texture2D, textures[1])
	p.ctx.ActiveTexture(texture2)
	p.ctx.BindTexture(texture2D, textures[2])
	p.ctx.ActiveTexture(texture3)
	p.ctx.BindTexture(texture2D, textures[3])
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithNVSemiplanar(textures [4]Texture, swapUV float32, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.nvSemiplanarProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.SetUniform1f(prog, "swapUV", swapUV)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, textures[0])
	p.ctx.ActiveTexture(texture1)
	p.ctx.BindTexture(texture2D, textures[1])
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithPackedYUV422(tex Texture, packingMode, texPackedWidth float32, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.packedYUV422Program

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.SetUniform1f(prog, "packingMode", packingMode)
	p.SetUniform1f(prog, "texPackedWidth", texPackedWidth)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, tex)
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithGrayscale(tex Texture, bitMax, hasAlpha float32, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.grayscaleProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.SetUniform1f(prog, "bitMax", bitMax)
	p.SetUniform1f(prog, "hasAlpha", hasAlpha)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, tex)
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithHiBitPlanar(textures [4]Texture, bitMax, hasAlpha float32, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.yuvPlanarHibitProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.SetUniform1f(prog, "bitMax", bitMax)
	p.SetUniform1f(prog, "hasAlpha", hasAlpha)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, textures[0])
	p.ctx.ActiveTexture(texture1)
	p.ctx.BindTexture(texture2D, textures[1])
	p.ctx.ActiveTexture(texture2)
	p.ctx.BindTexture(texture2D, textures[2])
	if hasAlpha > 0.5 {
		p.ctx.ActiveTexture(texture3)
		p.ctx.BindTexture(texture2D, textures[3])
	}
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithHiBitNV(textures [4]Texture, bitMax, swapUV float32, cm ColorMatrix, pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, 0)
	inner, _ := rectInnerCoords(size, pos, fill, 0)
	prog := p.nvSemiplanarHibitProgram

	p.ctx.UseProgram(prog.ref)
	p.setStreamProgramBase(prog, points, inner, insets, alpha)
	p.setColorMatrixUniforms(prog, cm)
	p.SetUniform1f(prog, "bitMax", bitMax)
	p.SetUniform1f(prog, "swapUV", swapUV)
	p.ctx.BlendFunc(one, oneMinusSrcAlpha)

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, textures[0])
	p.ctx.ActiveTexture(texture1)
	p.ctx.BindTexture(texture2D, textures[1])
	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.ctx.ActiveTexture(texture0)
	p.logError()
}

func (p *painter) drawQuadWithTexture(texture Texture, pos fyne.Position, size, frame fyne.Size,
	fill canvas.ImageFill, alpha, cornerRadius, pad float32,
) {
	points, insets := p.rectCoords(size, pos, frame, fill, 0, pad)
	inner, _ := rectInnerCoords(size, pos, fill, 0)

	p.ctx.UseProgram(p.program.ref)
	p.updateBuffer(p.program.buff, points)
	p.UpdateVertexArray(p.program, "vert", 3, 5, 0)
	p.UpdateVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	cornerRadius = fyne.Min(paint.GetMaximumRadius(size), cornerRadius)
	p.SetUniform1f(p.program, "cornerRadius", cornerRadius*p.pixScale)
	p.SetUniform2f(p.program, "size", inner.Width*p.pixScale, inner.Height*p.pixScale)
	p.SetUniform4f(p.program, "inset", insets[0], insets[1], insets[2], insets[3])

	p.SetUniform1f(p.program, "alpha", alpha)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) drawTextureWithDetails(o fyne.CanvasObject, creator func(canvasObject fyne.CanvasObject) Texture,
	pos fyne.Position, size, frame fyne.Size, fill canvas.ImageFill, alpha float32, pad float32,
) {
	texture, err := p.getTexture(o, creator)
	if err != nil {
		return
	}

	cornerRadius := float32(0)
	aspect := float32(0)
	if img, ok := o.(*canvas.Image); ok {
		aspect = img.Aspect()
		if aspect == 0 {
			aspect = 1 // fallback, should not occur - normally an image load error
		}
		if img.CornerRadius > 0 {
			cornerRadius = img.CornerRadius
		}
	}
	points, insets := p.rectCoords(size, pos, frame, fill, aspect, pad)
	inner, _ := rectInnerCoords(size, pos, fill, aspect)

	p.ctx.UseProgram(p.program.ref)
	p.updateBuffer(p.program.buff, points)
	p.UpdateVertexArray(p.program, "vert", 3, 5, 0)
	p.UpdateVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	// Set corner radius and texture size in pixels
	cornerRadius = fyne.Min(paint.GetMaximumRadius(size), cornerRadius)
	p.SetUniform1f(p.program, "cornerRadius", cornerRadius*p.pixScale)
	p.SetUniform2f(p.program, "size", inner.Width*p.pixScale, inner.Height*p.pixScale)
	p.SetUniform4f(p.program, "inset", insets[0], insets[1], insets[2], insets[3]) // texture coordinate insets (minX, minY, maxX, maxY)

	p.SetUniform1f(p.program, "alpha", alpha)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}

func (p *painter) lineCoords(pos, pos1, pos2 fyne.Position, lineWidth, feather float32, frame fyne.Size) ([]float32, float32, float32) {
	// Shift line coordinates so that they match the target position.
	xPosDiff := pos.X - fyne.Min(pos1.X, pos2.X)
	yPosDiff := pos.Y - fyne.Min(pos1.Y, pos2.Y)
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	pos2.X = roundToPixel(pos2.X+xPosDiff, p.pixScale)
	pos2.Y = roundToPixel(pos2.Y+yPosDiff, p.pixScale)

	if lineWidth <= 1 {
		offset := float32(0.5)                  // adjust location for lines < 1pt on regular display
		if lineWidth <= 0.5 && p.pixScale > 1 { // and for 1px drawing on HiDPI (width 0.5)
			offset = 0.25
		}
		if pos1.X == pos2.X {
			pos1.X -= offset
			pos2.X -= offset
		}
		if pos1.Y == pos2.Y {
			pos1.Y -= offset
			pos2.Y -= offset
		}
	}

	x1Pos := pos1.X / frame.Width
	x1 := -1 + x1Pos*2
	y1Pos := pos1.Y / frame.Height
	y1 := 1 - y1Pos*2
	x2Pos := pos2.X / frame.Width
	x2 := -1 + x2Pos*2
	y2Pos := pos2.Y / frame.Height
	y2 := 1 - y2Pos*2

	normalX := (pos2.Y - pos1.Y) / frame.Width
	normalY := (pos2.X - pos1.X) / frame.Height
	dirLength := float32(math.Sqrt(float64(normalX*normalX + normalY*normalY)))
	normalX /= dirLength
	normalY /= dirLength

	normalObjX := normalX * 0.5 * frame.Width
	normalObjY := normalY * 0.5 * frame.Height
	widthMultiplier := float32(math.Sqrt(float64(normalObjX*normalObjX + normalObjY*normalObjY)))
	halfWidth := (roundToPixel(lineWidth+feather, p.pixScale) * 0.5) / widthMultiplier
	featherWidth := feather / widthMultiplier

	return []float32{
		// coord x, y normal x, y
		x1, y1, normalX, normalY,
		x2, y2, normalX, normalY,
		x2, y2, -normalX, -normalY,
		x2, y2, -normalX, -normalY,
		x1, y1, normalX, normalY,
		x1, y1, -normalX, -normalY,
	}, halfWidth, featherWidth
}

// rectCoords calculates the openGL coordinate space of a rectangle
func (p *painter) rectCoords(size fyne.Size, pos fyne.Position, frame fyne.Size,
	fill canvas.ImageFill, aspect float32, pad float32,
) ([]float32, [4]float32) {
	size, pos = rectInnerCoords(size, pos, fill, aspect)
	size, pos = roundToPixelCoords(size, pos, p.pixScale)

	xPos := (pos.X - pad) / frame.Width
	x1 := -1 + xPos*2
	x2Pos := (pos.X + size.Width + pad) / frame.Width
	x2 := -1 + x2Pos*2

	yPos := (pos.Y - pad) / frame.Height
	y1 := 1 - yPos*2
	y2Pos := (pos.Y + size.Height + pad) / frame.Height
	y2 := 1 - y2Pos*2

	xInset := float32(0.0)
	yInset := float32(0.0)

	if fill == canvas.ImageFillCover {
		viewAspect := size.Width / size.Height

		if viewAspect > aspect {
			newHeight := size.Width / aspect
			heightPad := (newHeight - size.Height) / 2
			yInset = heightPad / newHeight
		} else if viewAspect < aspect {
			newWidth := size.Height * aspect
			widthPad := (newWidth - size.Width) / 2
			xInset = widthPad / newWidth
		}
	}

	insets := [4]float32{xInset, yInset, 1.0 - xInset, 1.0 - yInset}

	return []float32{
		// coord x, y, z texture x, y
		x1, y2, 0, insets[0], insets[3], // top left
		x1, y1, 0, insets[0], insets[1], // bottom left
		x2, y2, 0, insets[2], insets[3], // top right
		x2, y1, 0, insets[2], insets[1], // bottom right
	}, insets
}

func rectInnerCoords(size fyne.Size, pos fyne.Position, fill canvas.ImageFill, aspect float32) (fyne.Size, fyne.Position) {
	if fill == canvas.ImageFillContain || fill == canvas.ImageFillOriginal {
		// change pos and size accordingly

		viewAspect := size.Width / size.Height

		newWidth, newHeight := size.Width, size.Height
		widthPad, heightPad := float32(0), float32(0)
		if viewAspect > aspect {
			newWidth = size.Height * aspect
			widthPad = (size.Width - newWidth) / 2
		} else if viewAspect < aspect {
			newHeight = size.Width / aspect
			heightPad = (size.Height - newHeight) / 2
		}

		return fyne.NewSize(newWidth, newHeight), fyne.NewPos(pos.X+widthPad, pos.Y+heightPad)
	}

	return size, pos
}

func (p *painter) vecRectCoords(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size, aspect float32) ([4]float32, []float32) {
	xPad, yPad := float32(0), float32(0)

	if aspect != 0 {
		inner := rect.Size()
		frameAspect := inner.Width / inner.Height

		if frameAspect > aspect {
			newWidth := inner.Height * aspect
			xPad = (inner.Width - newWidth) / 2
		} else if frameAspect < aspect {
			newHeight := inner.Width / aspect
			yPad = (inner.Height - newHeight) / 2
		}
	}

	return p.vecRectCoordsWithPad(pos, rect, frame, xPad, yPad)
}

func (p *painter) vecRectCoordsWithPad(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size, xPad, yPad float32) ([4]float32, []float32) {
	size := rect.Size()
	pos1 := rect.Position()

	xPosDiff := pos.X - pos1.X + xPad
	yPosDiff := pos.Y - pos1.Y + yPad
	pos1.X = roundToPixel(pos1.X+xPosDiff, p.pixScale)
	pos1.Y = roundToPixel(pos1.Y+yPosDiff, p.pixScale)
	size.Width = roundToPixel(size.Width-2*xPad, p.pixScale)
	size.Height = roundToPixel(size.Height-2*yPad, p.pixScale)

	// without edge softness adjustment the rectangle has cropped edges
	edgeSoftnessScaled := roundToPixel(p.edgeSoftness()*p.pixScale, 1.0)
	x1Pos := pos1.X
	x1Norm := -1 + (x1Pos-edgeSoftnessScaled)*2/frame.Width
	x2Pos := pos1.X + size.Width
	x2Norm := -1 + (x2Pos+edgeSoftnessScaled)*2/frame.Width
	y1Pos := pos1.Y
	y1Norm := 1 - (y1Pos-edgeSoftnessScaled)*2/frame.Height
	y2Pos := pos1.Y + size.Height
	y2Norm := 1 - (y2Pos+edgeSoftnessScaled)*2/frame.Height

	// output a norm for the fill and the vert is unused, but we pass 0 to avoid optimisation issues
	coords := []float32{
		0, 0, x1Norm, y1Norm, // first triangle
		0, 0, x2Norm, y1Norm, // second triangle
		0, 0, x1Norm, y2Norm,
		0, 0, x2Norm, y2Norm,
	}

	return [4]float32{x1Pos, y1Pos, x2Pos, y2Pos}, coords
}

func (p *painter) vecSquareCoords(pos fyne.Position, rect fyne.CanvasObject, frame fyne.Size) ([4]float32, []float32) {
	return p.vecRectCoordsWithPad(pos, rect, frame, 0, 0)
}

func roundToPixel(v float32, pixScale float32) float32 {
	if pixScale == 1.0 {
		return float32(math.Round(float64(v)))
	}

	return float32(math.Round(float64(v*pixScale))) / pixScale
}

func roundToPixelCoords(size fyne.Size, pos fyne.Position, pixScale float32) (fyne.Size, fyne.Position) {
	end := pos.Add(size)
	end.X = roundToPixel(end.X, pixScale)
	end.Y = roundToPixel(end.Y, pixScale)
	pos.X = roundToPixel(pos.X, pixScale)
	pos.Y = roundToPixel(pos.Y, pixScale)
	size.Width = end.X - pos.X
	size.Height = end.Y - pos.Y

	return size, pos
}

// Returns FragmentColor(red,green,blue,alpha) from fyne.Color
func getFragmentColor(col color.Color) (float32, float32, float32, float32) {
	if col == nil {
		return 0, 0, 0, 0
	}
	r, g, b, a := col.RGBA()
	if a == 0 {
		return 0, 0, 0, 0
	}
	alpha := float32(a)
	return float32(r) / alpha, float32(g) / alpha, float32(b) / alpha, alpha / 0xffff
}

func (p *painter) scaleFrameSize(frame fyne.Size) (float32, float32) {
	frameWidthScaled := roundToPixel(frame.Width*p.pixScale, 1.0)
	frameHeightScaled := roundToPixel(frame.Height*p.pixScale, 1.0)
	return frameWidthScaled, frameHeightScaled
}

// Returns scaled RectCoords(x1,x2,y1,y2) in same order
func (p *painter) scaleRectCoords(x1, x2, y1, y2 float32) (float32, float32, float32, float32) {
	x1Scaled := roundToPixel(x1*p.pixScale, 1.0)
	x2Scaled := roundToPixel(x2*p.pixScale, 1.0)
	y1Scaled := roundToPixel(y1*p.pixScale, 1.0)
	y2Scaled := roundToPixel(y2*p.pixScale, 1.0)
	return x1Scaled, x2Scaled, y1Scaled, y2Scaled
}
