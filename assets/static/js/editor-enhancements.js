document.addEventListener('DOMContentLoaded', () => {
    const body = document.body;
    const markdownPane = document.getElementById('markdown-pane');

    const staticZenButton = document.getElementById('zen-mode-btn-static');
    const floatingZenButton = document.getElementById('zen-mode-btn-floating');
    const floatingDarkModeButton = document.getElementById('dark-mode-btn-floating');

    function enterZenMode() {
        body.classList.add('zen-mode-active');
    }

    function exitZenMode() {
        body.classList.remove('zen-mode-active');
        if (markdownPane.classList.contains('dark-mode-active')) {
            markdownPane.classList.remove('dark-mode-active');
            floatingDarkModeButton.textContent = '🌙';
            floatingDarkModeButton.title = 'Dark mode (alt+d)';
        }
    }

    function toggleDarkMode() {
        markdownPane.classList.toggle('dark-mode-active');
        const isDark = markdownPane.classList.contains('dark-mode-active');
        floatingDarkModeButton.textContent = isDark ? '☀️' : '🌙';
    }

    if (staticZenButton) {
        staticZenButton.addEventListener('click', enterZenMode);
    }

    if (floatingZenButton) {
        floatingZenButton.addEventListener('click', exitZenMode);
    }

    if (floatingDarkModeButton) {
        floatingDarkModeButton.addEventListener('click', toggleDarkMode);
    }

    document.addEventListener('keydown', (event) => {
        const isZen = body.classList.contains('zen-mode-active');

        if (event.altKey && event.key.toLowerCase() === 'z') {
            event.preventDefault();
            if (isZen) {
                exitZenMode();
            } else {
                enterZenMode();
            }
        }

        if (event.altKey && event.key.toLowerCase() === 'd') {
            if (isZen) {
                event.preventDefault();
                toggleDarkMode();
            }
        }
    });
});
