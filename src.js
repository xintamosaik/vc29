
import AnnotationScript from './pages/annotate.js';

// Make the function globally available for the templ onmouseup handler (neeed for esbuild to not treeshake it)
window.handleAnnotateMouseUp = AnnotationScript;
(function initNavigation() {
    const navigation = window.main_navigation;
    const links = navigation.querySelectorAll('a');
    function setActive(event) {
        if (event.target.tagName !== 'A') {
            return;
        }
        links.forEach(link => {
            if (link === event.target) {
                link.classList.add('active');
            } else {
                link.classList.remove('active');
            }
        });
    }
    navigation.addEventListener('click', setActive);
})();   
