(function initNavigation() {
    const navigation = window.main_navigation;
    const links = navigation.querySelectorAll('a');

    function navigate(event) {
        const clicked = event.target
        if (clicked.tagName !== 'A') {
            return;
        }

        const sections = document.querySelectorAll('section.main');
        sections.forEach(section => {
            if (section.id === clicked.dataset.reference) {
                section.classList.remove('hidden');
            } else {
                section.classList.add('hidden');
            }
        });

        links.forEach(link => {
            if (link === clicked) {
                link.classList.add('nav-active')
            } else {
                link.classList.remove('nav-active');
            }
        });
    }
    navigation.addEventListener('click', navigate);
})();   