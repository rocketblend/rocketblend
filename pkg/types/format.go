package types

// RenderFormat represents the available formats for rendering.
type RenderFormat string

const (
	RenderFormatTGA                 RenderFormat = "TGA"
	RenderFormatRAWTGA              RenderFormat = "RAWTGA"
	RenderFormatJPEG                RenderFormat = "JPEG"
	RenderFormatIRIS                RenderFormat = "IRIS"
	RenderFormatAVIRAW              RenderFormat = "AVIRAW"
	RenderFormatAVIJPEG             RenderFormat = "AVIJPEG"
	RenderFormatPNG                 RenderFormat = "PNG"
	RenderFormatBMP                 RenderFormat = "BMP"
	RenderFormatHDR                 RenderFormat = "HDR"
	RenderFormatTIFF                RenderFormat = "TIFF"
	RenderFormatOPEN_EXR            RenderFormat = "OPEN_EXR"
	RenderFormatOPEN_EXR_MULTILAYER RenderFormat = "OPEN_EXR_MULTILAYER"
	RenderFormatFFMPEG              RenderFormat = "FFMPEG"
	RenderFormatCINEON              RenderFormat = "CINEON"
	RenderFormatDPX                 RenderFormat = "DPX"
	RenderFormatJP2                 RenderFormat = "JP2"
	RenderFormatWEBP                RenderFormat = "WEBP"
)
