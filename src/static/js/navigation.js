// Gerenciador de navegação
class Navigation {
    static init() {
        this.highlightActiveLink();
        this.bindMobileMenu();
    }

    static highlightActiveLink() {
        const currentPath = window.location.pathname;
        const navLinks = document.querySelectorAll('.nav-link');
        
        navLinks.forEach(link => {
            const href = link.getAttribute('href');
            if (href === currentPath || (currentPath === '/' && href === '/')) {
                link.classList.add('active');
            } else {
                link.classList.remove('active');
            }
        });
    }

    static bindMobileMenu() {
        // Adicionar toggle para menu mobile se necessário
        const mobileMenuBtn = document.getElementById('mobile-menu-btn');
        if (mobileMenuBtn) {
            mobileMenuBtn.addEventListener('click', () => {
                const navLinks = document.querySelector('.nav-links');
                navLinks.classList.toggle('active');
            });
        }
    }

    static navigateTo(url) {
        window.location.href = url;
    }

    static reload() {
        window.location.reload();
    }
}