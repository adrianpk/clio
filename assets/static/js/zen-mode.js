document.addEventListener('DOMContentLoaded', () => {
    const zenButton = document.getElementById('zen-mode-btn');
    const editorTextarea = document.getElementById('body');

    if (zenButton && editorTextarea) {
        zenButton.addEventListener('click', () => {
            document.body.classList.toggle('zen-mode-active');

            if (document.body.classList.contains('zen-mode-active')) {
                zenButton.textContent = '💻';
            } else {
                zenButton.textContent = '🧘';
            }
        });
    }
});
