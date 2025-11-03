// JavaScript para o app dash
console.log('App dash carregado!');

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM carregado para o app dash');
    
    // Adicione a lógica JavaScript do seu app aqui
    
    // Exemplo: interação básica
    const welcomeCard = document.querySelector('.welcome-card');
    if (welcomeCard) {
        welcomeCard.addEventListener('click', function() {
            this.style.transform = 'scale(0.98)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 150);
        });
    }
});
