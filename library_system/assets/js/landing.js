document.addEventListener('DOMContentLoaded', () => {
    const buttons = document.querySelectorAll('.main-btn');

    buttons.forEach(button => {
        button.addEventListener('mouseenter', (e) => {
            e.target.style.transform = 'scale(1.05)';
        });

        button.addEventListener('mouseleave', (e) => {
            e.target.style.transform = 'scale(1)';
        });
    });
});
