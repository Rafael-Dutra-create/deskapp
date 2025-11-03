// About Page JavaScript
class AboutPage {
    static init() {
        console.log('📄 AboutPage inicializado');
        
        this.bindEvents();
        this.loadAdditionalInfo();
        this.initAnimations();
        this.initInteractiveElements(); // Agora este método existe
    }

    static bindEvents() {
        // Botão de detalhes técnicos
        const detailsBtn = document.querySelector('.btn-secondary');
        if (detailsBtn) {
            detailsBtn.addEventListener('click', this.toggleTechnicalDetails);
        }

        // Cards de features - hover effects
        this.initFeatureCards();
        
        // Team cards - interações
        this.initTeamCards();
        
        // Copy email functionality
        this.initEmailCopy();
        
        // Statistics animation
        this.initStatsAnimation();
    }

    // NOVO MÉTODO: Inicializar elementos interativos
    static initInteractiveElements() {
        console.log('🎮 Inicializando elementos interativos');
        
        // Tooltips para elementos com data-tooltip
        this.initTooltips();
        
        // Accordion para seções técnicas (se houver)
        this.initAccordions();
        
        // Lazy loading para imagens (se adicionadas futuramente)
        this.initLazyLoading();
        
        // Botões de ação rápida
        this.initQuickActions();
    }

    // NOVO: Inicializar tooltips
    static initTooltips() {
        const tooltipElements = document.querySelectorAll('[data-tooltip]');
        
        tooltipElements.forEach(element => {
            element.addEventListener('mouseenter', this.showTooltip);
            element.addEventListener('mouseleave', this.hideTooltip);
            element.addEventListener('focus', this.showTooltip);
            element.addEventListener('blur', this.hideTooltip);
        });
    }

    static showTooltip(e) {
        const element = e.target;
        const tooltipText = element.getAttribute('data-tooltip');
        
        // Remover tooltip existente
        this.hideTooltip();
        
        // Criar tooltip
        const tooltip = document.createElement('div');
        tooltip.className = 'custom-tooltip';
        tooltip.textContent = tooltipText;
        tooltip.id = 'current-tooltip';
        
        document.body.appendChild(tooltip);
        
        // Posicionar tooltip
        const rect = element.getBoundingClientRect();
        tooltip.style.left = `${rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2)}px`;
        tooltip.style.top = `${rect.top - tooltip.offsetHeight - 10}px`;
    }

    static hideTooltip() {
        const existingTooltip = document.getElementById('current-tooltip');
        if (existingTooltip) {
            existingTooltip.remove();
        }
    }

    // NOVO: Inicializar accordions
    static initAccordions() {
        const accordionHeaders = document.querySelectorAll('.accordion-header');
        
        accordionHeaders.forEach(header => {
            header.addEventListener('click', () => {
                const content = header.nextElementSibling;
                const isOpen = content.style.maxHeight;
                
                // Fechar todos os accordions
                document.querySelectorAll('.accordion-content').forEach(acc => {
                    acc.style.maxHeight = null;
                    acc.previousElementSibling.classList.remove('active');
                });
                
                // Abrir o clicado se não estava aberto
                if (!isOpen) {
                    content.style.maxHeight = content.scrollHeight + "px";
                    header.classList.add('active');
                }
            });
        });
    }

    // NOVO: Inicializar lazy loading
    static initLazyLoading() {
        const lazyImages = document.querySelectorAll('img[data-src]');
        
        if ('IntersectionObserver' in window) {
            const imageObserver = new IntersectionObserver((entries, observer) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        const img = entry.target;
                        img.src = img.getAttribute('data-src');
                        img.classList.remove('lazy');
                        imageObserver.unobserve(img);
                    }
                });
            });

            lazyImages.forEach(img => imageObserver.observe(img));
        } else {
            // Fallback para browsers antigos
            lazyImages.forEach(img => {
                img.src = img.getAttribute('data-src');
            });
        }
    }

    // NOVO: Inicializar ações rápidas
    static initQuickActions() {
        // Botão de compartilhar
        const shareBtn = document.querySelector('[data-action="share"]');
        if (shareBtn) {
            shareBtn.addEventListener('click', this.sharePage);
        }
        
        // Botão de imprimir
        const printBtn = document.querySelector('[data-action="print"]');
        if (printBtn) {
            printBtn.addEventListener('click', this.printPage);
        }
        
        // Botão de toggle dark mode
        const themeBtn = document.querySelector('[data-action="toggle-theme"]');
        if (themeBtn) {
            themeBtn.addEventListener('click', this.toggleDarkMode);
        }
    }

    static toggleTechnicalDetails() {
        const url = new URL(window.location);
        const currentDetailed = url.searchParams.get('detailed') === 'true';
        
        if (currentDetailed) {
            url.searchParams.delete('detailed');
        } else {
            url.searchParams.set('detailed', 'true');
        }
        
        // Mostrar loading
        AppUtils.showNotification('Carregando detalhes...', 'info');
        
        // Navegar para a mesma página com parâmetro diferente
        setTimeout(() => {
            window.location.href = url.toString();
        }, 500);
    }

    static initFeatureCards() {
        const featureCards = document.querySelectorAll('.feature-card');
        
        featureCards.forEach(card => {
            // Efeito de tilt suave
            card.addEventListener('mousemove', this.handleCardTilt);
            card.addEventListener('mouseleave', this.resetCardTilt);
            
            // Click para expandir
            card.addEventListener('click', (e) => {
                if (!e.target.closest('a')) {
                    this.toggleFeatureDescription(card);
                }
            });
        });
    }

    static handleCardTilt(e) {
        const card = e.currentTarget;
        const rect = card.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        const centerX = rect.width / 2;
        const centerY = rect.height / 2;
        
        const rotateY = (x - centerX) / 25;
        const rotateX = (centerY - y) / 25;
        
        card.style.transform = `perspective(1000px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) scale3d(1.02, 1.02, 1.02)`;
    }

    static resetCardTilt(e) {
        const card = e.currentTarget;
        card.style.transform = 'perspective(1000px) rotateX(0) rotateY(0) scale3d(1, 1, 1)';
    }

    static toggleFeatureDescription(card) {
        const description = card.querySelector('p');
        const isExpanded = card.classList.contains('expanded');
        
        if (isExpanded) {
            description.style.maxHeight = '3em';
            card.classList.remove('expanded');
        } else {
            description.style.maxHeight = `${description.scrollHeight}px`;
            card.classList.add('expanded');
        }
    }

    static initTeamCards() {
        const teamCards = document.querySelectorAll('.team-card');
        
        teamCards.forEach(card => {
            card.addEventListener('click', () => {
                this.showTeamMemberModal(card);
            });
            
            // Efeito de hover com delay
            card.addEventListener('mouseenter', () => {
                setTimeout(() => {
                    if (card.matches(':hover')) {
                        card.classList.add('hover-active');
                    }
                }, 100);
            });
            
            card.addEventListener('mouseleave', () => {
                card.classList.remove('hover-active');
            });
        });
    }

    static showTeamMemberModal(card) {
        const name = card.querySelector('h3').textContent;
        const role = card.querySelector('.team-role').textContent;
        const email = card.querySelector('.team-email').textContent;
        const expertise = card.querySelector('.team-expertise').textContent;
        const avatar = card.querySelector('.team-avatar').textContent;
        
        const modalHTML = `
            <div class="modal-overlay" onclick="AboutPage.closeModal()">
                <div class="modal-content" onclick="event.stopPropagation()">
                    <button class="modal-close" onclick="AboutPage.closeModal()">×</button>
                    
                    <div class="modal-header">
                        <div class="modal-avatar">${avatar}</div>
                        <div class="modal-title">
                            <h2>${name}</h2>
                            <p class="modal-role">${role}</p>
                        </div>
                    </div>
                    
                    <div class="modal-body">
                        <div class="contact-info">
                            <h4>Contato</h4>
                            <p class="contact-email" onclick="AboutPage.copyEmail('${email}')">
                                📧 ${email}
                                <span class="copy-hint">(clique para copiar)</span>
                            </p>
                        </div>
                        
                        <div class="expertise-info">
                            <h4>Especialidades</h4>
                            <p>${expertise}</p>
                        </div>
                        
                        <div class="action-buttons">
                            <button class="btn btn-outline" onclick="AboutPage.scheduleMeeting('${name}')">
                                📅 Agendar Reunião
                            </button>
                            <button class="btn btn-outline" onclick="AboutPage.viewProfile('${name}')">
                                👤 Ver Perfil Completo
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;
        
        document.body.insertAdjacentHTML('beforeend', modalHTML);
        this.showModal();
    }

    static showModal() {
        const modal = document.querySelector('.modal-overlay');
        if (modal) {
            modal.style.display = 'flex';
            setTimeout(() => {
                modal.classList.add('active');
            }, 10);
            
            // Adicionar evento para fechar com ESC
            document.addEventListener('keydown', this.handleEscKey);
        }
    }

    static closeModal() {
        const modal = document.querySelector('.modal-overlay');
        if (modal) {
            modal.classList.remove('active');
            setTimeout(() => {
                modal.remove();
                document.removeEventListener('keydown', this.handleEscKey);
            }, 300);
        }
    }

    static handleEscKey = (e) => {
        if (e.key === 'Escape') {
            this.closeModal();
        }
    }

    static initEmailCopy() {
        const emailElements = document.querySelectorAll('.team-email');
        emailElements.forEach(element => {
            element.style.cursor = 'pointer';
            element.title = 'Clique para copiar email';
            
            element.addEventListener('click', (e) => {
                e.stopPropagation();
                this.copyEmail(element.textContent);
            });
        });
    }

    static copyEmail(email) {
        navigator.clipboard.writeText(email).then(() => {
            AppUtils.showNotification('Email copiado para a área de transferência!', 'success');
        }).catch(err => {
            console.error('Falha ao copiar email:', err);
            AppUtils.showNotification('Erro ao copiar email', 'error');
        });
    }

    static scheduleMeeting(memberName) {
        AppUtils.showNotification(`Reunião com ${memberName} agendada! (simulação)`, 'success');
        this.closeModal();
        
        setTimeout(() => {
            if (confirm(`Deseja abrir o calendário para agendar com ${memberName}?`)) {
                window.open('https://calendar.google.com', '_blank');
            }
        }, 1000);
    }

    static viewProfile(memberName) {
        AppUtils.showNotification(`Perfil de ${memberName} carregado! (simulação)`, 'info');
        this.closeModal();
    }

    static initStatsAnimation() {
        const statsSection = document.querySelector('.stats-section');
        if (!statsSection) return;

        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    this.animateStats();
                    observer.unobserve(entry.target);
                }
            });
        }, { threshold: 0.3 });

        observer.observe(statsSection);
    }

    static animateStats() {
        const statNumbers = document.querySelectorAll('.stat-details strong');
        
        statNumbers.forEach((element, index) => {
            const originalText = element.textContent;
            const value = parseInt(originalText.replace(/\D/g, ''));
            
            if (!isNaN(value)) {
                this.animateValue(element, 0, value, 1500 + (index * 200));
            }
        });
    }

    static animateValue(element, start, end, duration) {
        const startTime = performance.now();
        const unit = element.textContent.replace(/[0-9]/g, '') || '';
        
        function updateValue(currentTime) {
            const elapsed = currentTime - startTime;
            const progress = Math.min(elapsed / duration, 1);
            
            const easeOutQuart = 1 - Math.pow(1 - progress, 4);
            const currentValue = Math.floor(start + (end - start) * easeOutQuart);
            
            element.textContent = currentValue + unit;
            
            if (progress < 1) {
                requestAnimationFrame(updateValue);
            }
        }
        
        requestAnimationFrame(updateValue);
    }

    static initAnimations() {
        const animatedElements = document.querySelectorAll('.feature-card, .tech-category, .team-card');
        
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('animate-in');
                    observer.unobserve(entry.target);
                }
            });
        }, { threshold: 0.1 });

        animatedElements.forEach(element => {
            observer.observe(element);
        });
    }

    static loadAdditionalInfo() {
        if (window.location.search.includes('detailed=true')) {
            this.loadDetailedStats();
        }
    }

    static loadDetailedStats() {
        fetch('/api/about')
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    this.updateDetailedInfo(data.data);
                }
            })
            .catch(error => {
                console.error('Erro ao carregar estatísticas detalhadas:', error);
            });
    }

    static updateDetailedInfo(data) {
        console.log('Informações detalhadas carregadas:', data);
        
        const serverTimeElement = document.createElement('div');
        serverTimeElement.className = 'server-time';
        serverTimeElement.innerHTML = `
            <small>Última atualização: ${new Date(data.server_time).toLocaleString('pt-BR')}</small>
        `;
        
        const ctaSection = document.querySelector('.cta-section');
        if (ctaSection) {
            ctaSection.appendChild(serverTimeElement);
        }
    }

    // Métodos de ações rápidas
    static sharePage() {
        const title = document.title;
        const url = window.location.href;
        
        if (navigator.share) {
            navigator.share({
                title: title,
                url: url
            }).then(() => {
                AppUtils.showNotification('Página compartilhada com sucesso!', 'success');
            }).catch(err => {
                console.log('Erro ao compartilhar:', err);
            });
        } else {
            this.copyEmail(url);
            AppUtils.showNotification('URL copiada para a área de transferência!', 'success');
        }
    }

    static printPage() {
        window.print();
    }

    static toggleDarkMode() {
        const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        ThemeManager.setTheme(newTheme);
        AppUtils.showNotification(`Modo ${newTheme === 'dark' ? 'escuro' : 'claro'} ativado`, 'success');
    }
}

// Inicialização quando o script carregar
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => AboutPage.init());
} else {
    AboutPage.init();
}

// Funções globais para acesso via HTML
window.showTechnicalDetails = () => AboutPage.toggleTechnicalDetails();
window.AboutPage = AboutPage;