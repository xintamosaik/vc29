
import AnnotationScript from './internal/pages/annotate.js';

// Make the function globally available for the templ onmouseup handler (neeed for esbuild to not treeshake it)
window.handleAnnotateMouseUp = AnnotationScript;

