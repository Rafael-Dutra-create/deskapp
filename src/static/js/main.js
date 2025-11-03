// Aplicação principal
class MeuApp {
    static init() {
        console.log('🎯 MeuApp inicializando...');
        
        // Inicializar componentes
        this.initComponents();
        this.bindGlobalEvents();
        this.checkPerformance();
    }

    static initComponents() {
        // Inicializar tooltips
        this.initTooltips();
        
        // Inicializar loaders
        this.initLoaders();
        
        // Inicializar modais
        this.initModals();
    }

    static initTooltips() {
        const tooltips = document.querySelectorAll('[data-tooltip]');
        tooltips.forEach(element => {
            element.addEventListener('mouseenter', this.showTooltip);
            element.addEventListener('mouseleave', this.hideTooltip);
        });
    }

    static showTooltip(e) {
        const text = this.getAttribute('data-tooltip');
        const tooltip = document.createElement('div');
        tooltip.className = 'tooltip';
        tooltip.textContent = text;
        document.body.appendChild(tooltip);
        
        const rect = this.getBoundingClientRect();
        tooltip.style.left = rect.left + 'px';
        tooltip.style.top = (rect.top - tooltip.offsetHeight - 5) + 'px';
    }

    static hideTooltip() {
        const tooltip = document.querySelector('.tooltip');
        if (tooltip) {
            tooltip.remove();
        }
    }

    static initLoaders() {
        // Interceptar links para mostrar loading
        document.addEventListener('click', (e) => {
            const link = e.target.closest('a');
            if (link && link.href && !link.target) {
                this.showLoader();
            }
        });
    }

    static showLoader() {
        const loader = document.createElement('div');
        loader.id = 'global-loader';
        loader.innerHTML = `
            <div class="loader-spinner"></div>
            <p>Carregando...</p>
        `;
        document.body.appendChild(loader);
    }

    static hideLoader() {
        const loader = document.getElementById('global-loader');
        if (loader) {
            loader.remove();
        }
    }

    static initModals() {
        // Fechar modal com ESC
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                this.closeAllModals();
            }
        });
    }

    static closeAllModals() {
        const modals = document.querySelectorAll('.modal');
        modals.forEach(modal => {
            modal.style.display = 'none';
        });
    }

    static bindGlobalEvents() {
        // Esconder loader quando a página carregar
        window.addEventListener('load', () => {
            this.hideLoader();
        });

        // Tratar erros globais
        window.addEventListener('error', (e) => {
            console.error('Erro global:', e.error);
            AppUtils.showNotification('Ocorreu um erro inesperado', 'error');
        });

        // Online/Offline detection
        window.addEventListener('online', () => {
            AppUtils.showNotification('Conexão restaurada', 'success');
        });

        window.addEventListener('offline', () => {
            AppUtils.showNotification('Conexão perdida', 'warning');
        });
    }

    static checkPerformance() {
        // Performance monitoring
        window.addEventListener('load', () => {
            const perfData = performance.timing;
            const loadTime = perfData.loadEventEnd - perfData.navigationStart;
            console.log(`📊 Page loaded in ${loadTime}ms`);
            
            if (loadTime > 3000) {
                console.warn('⚠️  Page load time is slow');
            }
        });
    }
}

// Inicializar app quando o DOM estiver pronto
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => MeuApp.init());
} else {
    MeuApp.init();
}