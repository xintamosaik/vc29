(function initNavigation() {
    const navigation = window.main_navigation;
    const links = navigation.querySelectorAll('a');

    function navigate(event) {
        const sections = document.querySelectorAll('section.main');
        if (event.target.tagName !== 'A') {
            return;
        }

        const target = event.target.dataset.reference;
        console.log(`Navigating to: ${target}`);

        sections.forEach(section => {
            console.log(`Checking section: ${section.id}`);
            if (section.id === target) {
                console.log(`Showing section: ${section.id}`);
                section.classList.remove('hidden');
            } else {
                console.log(`Hiding section: ${section.id}`);
                section.classList.add('hidden');
            }
        });



        links.forEach(link => {
            if (link === event.target) {
                link.classList.add('nav-active')
            } else {
                link.classList.remove('nav-active');
            }
        });


    }
    navigation.addEventListener('click', navigate);
})();   