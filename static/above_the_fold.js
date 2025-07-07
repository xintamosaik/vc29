(function initNavigation() {
    const navigation = window.main_navigation;
    const links = navigation.querySelectorAll('a');
    function setActive(event) {
        if (event.target.tagName !== 'A') {
            return;
        }
        links.forEach(link => {
            if (link === event.target) {
                link.classList.add('nav-active')
            } else {
                link.classList.remove('nav-active');
            }
        });
    }
    navigation.addEventListener('click', setActive);
})();   