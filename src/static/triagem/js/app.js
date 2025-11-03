// JavaScript para o app triagem
console.log('App triagem (Gin) carregado!');

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM carregado para o app triagem');
    
    // Adicione a lógica JavaScript do seu app aqui
    
    // Exemplo: interação básica
    const welcomeCard = document.querySelector('.welcome-card');
    if (welcomeCard) {
        welcomeCard.addEventListener('mouseenter', function() {
            this.style.boxShadow = '0 8px 12px rgba(0, 0, 0, 0.15)';
        });
		welcomeCard.addEventListener('mouseleave', function() {
            this.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
        });
    }
});
