// Inicializar AOS
AOS.init({
    duration: 600,
    once: true
});

let currentStep = 1;
const totalSteps = 3;

function nextStep(step) {
    if (validateCurrentStep()) {
        showStep(step);
    }
}

function prevStep(step) {
    showStep(step);
}

function showStep(step) {
    // Esconder todos os conteúdos
    document.querySelectorAll('.step-content').forEach(content => {
        content.classList.add('d-none');
    });
    
    // Mostrar conteúdo atual
    document.getElementById(`stepContent${step}`).classList.remove('d-none');
    
    // Atualizar indicadores
    updateStepIndicators(step);
    
    currentStep = step;
}

function updateStepIndicators(activeStep) {
    for (let i = 1; i <= totalSteps; i++) {
        const stepItem = document.getElementById(`stepItem${i}`);
        if (!stepItem) continue;

        stepItem.classList.remove('active', 'completed');

        if (i < activeStep) {
            stepItem.classList.add('completed');
        } else if (i === activeStep) {
            stepItem.classList.add('active');
        }
    }
}

function validateCurrentStep() {
    const currentContent = document.getElementById(`stepContent${currentStep}`);
    const inputs = currentContent.querySelectorAll('input[required], select[required]');
    let isValid = true;
    
    inputs.forEach(input => {
        if (!input.checkValidity()) {
            isValid = false;
            input.classList.add('is-invalid');
        } else {
            input.classList.remove('is-invalid');
        }
    });
    
    return isValid;
}

// Funções para manipulação de CPF/CNPJ
function formatCPF(cpf) {
    cpf = cpf.replace(/\D/g, '');
    cpf = cpf.replace(/([\d]{3})([\d]{3})([\d]{3})([\d]{2})/, "$1.$2.$3-$4");
    return cpf;
}

function formatCNPJ(cnpj) {
    cnpj = cnpj.replace(/\D/g, '');
    cnpj = cnpj.replace(/([\d]{2})([\d]{3})([\d]{3})([\d]{4})([\d]{2})/, "$1.$2.$3/$4-$5");
    return cnpj;
}

function isValidCPF(cpf) {
    cpf = cpf.replace(/\D/g, '');
    if (cpf.length !== 11 || /^(\d)\1{10}$/.test(cpf)) return false;
    
    let sum = 0;
    let remainder;
    
    for (let i = 1; i <= 9; i++) {
        sum += parseInt(cpf.substring(i-1, i)) * (11 - i);
    }
    
    remainder = (sum * 10) % 11;
    if (remainder === 10 || remainder === 11) remainder = 0;
    if (remainder !== parseInt(cpf.substring(9, 10))) return false;
    
    sum = 0;
    for (let i = 1; i <= 10; i++) {
        sum += parseInt(cpf.substring(i-1, i)) * (12 - i);
    }
    
    remainder = (sum * 10) % 11;
    if (remainder === 10 || remainder === 11) remainder = 0;
    if (remainder !== parseInt(cpf.substring(10, 11))) return false;
    
    return true;
}

function isValidCNPJ(cnpj) {
    cnpj = cnpj.replace(/\D/g, '');
    if (cnpj.length !== 14 || /^(\d)\1{13}$/.test(cnpj)) return false;
    
    let size = cnpj.length - 2;
    let numbers = cnpj.substring(0, size);
    let digits = cnpj.substring(size);
    let sum = 0;
    let pos = size - 7;
    
    for (let i = size; i >= 1; i--) {
        sum += numbers.charAt(size - i) * pos--;
        if (pos < 2) pos = 9;
    }
    
    let result = sum % 11 < 2 ? 0 : 11 - sum % 11;
    if (result !== parseInt(digits.charAt(0))) return false;
    
    size = size + 1;
    numbers = cnpj.substring(0, size);
    sum = 0;
    pos = size - 7;
    
    for (let i = size; i >= 1; i--) {
        sum += numbers.charAt(size - i) * pos--;
        if (pos < 2) pos = 9;
    }
    
    result = sum % 11 < 2 ? 0 : 11 - sum % 11;
    if (result !== parseInt(digits.charAt(1))) return false;
    
    return true;
}

// Validação de força da senha
document.getElementById('senha').addEventListener('input', function() {
    const password = this.value;
    const strengthBar = document.getElementById('passwordStrength');
    
    let strength = 0;
    if (password.length >= 8) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[^A-Za-z0-9]/.test(password)) strength++;
    
    strengthBar.className = 'password-strength ';
    if (strength <= 1) {
        strengthBar.classList.add('strength-weak');
    } else if (strength <= 2) {
        strengthBar.classList.add('strength-medium');
    } else {
        strengthBar.classList.add('strength-strong');
    }
});

// Validação de confirmação de senha
document.getElementById('confirmarSenha').addEventListener('input', function() {
    const senha = document.getElementById('senha').value;
    const confirmar = this.value;
    
    if (senha !== confirmar) {
        this.setCustomValidity('As senhas não coincidem');
    } else {
        this.setCustomValidity('');
    }
});

// Máscara para telefone
document.getElementById('telefone').addEventListener('input', function() {
    let value = this.value.replace(/\D/g, '');
    value = value.replace(/(\d{2})(\d)/, '($1) $2');
    value = value.replace(/(\d{5})(\d)/, '$1-$2');
    this.value = value;
});

// Função para consultar CNPJ na API
async function consultarCNPJ(cnpj) {
    try {
        const response = await fetch(`https://brasilapi.com.br/api/cnpj/v1/${cnpj}`);
        if (!response.ok) {
            throw new Error('Erro ao consultar CNPJ');
        }
        return await response.json();
    } catch (error) {
        console.error('Erro na consulta de CNPJ:', error);
        return null;
    }
}

// Máscara para CEP
document.getElementById('cep').addEventListener('input', function() {
    let value = this.value.replace(/\D/g, '');
    value = value.replace(/(\d{5})(\d)/, '$1-$2');
    this.value = value;
});

// Buscar endereço por CEP
document.getElementById('cep').addEventListener('blur', function() {
    const cep = this.value.replace(/\D/g, '');
    if (cep.length === 8) {
        fetch(`https://viacep.com.br/ws/${cep}/json/`)
            .then(response => response.json())
            .then(data => {
                if (!data.erro) {
                    document.getElementById('endereco').value = data.logradouro;
                    document.getElementById('bairro').value = data.bairro;
                    document.getElementById('cidade').value = data.localidade;
                    document.getElementById('estado').value = data.uf;
                }
            })
            .catch(error => console.error('Erro ao buscar CEP:', error));
    }
});

// Evento para alternar entre CPF e CNPJ
document.addEventListener('DOMContentLoaded', function() {
    const tipoPessoaSelect = document.getElementById('tipoPessoa');
    const documentoInput = document.getElementById('documento');
    const documentoLabel = document.getElementById('documentoLabel');
    const nomeLabel = document.getElementById('nomeLabel');
    const nomeFantasiaContainer = document.getElementById('nomeFantasiaContainer');
    const pessoaFisicaElements = document.querySelectorAll('.pessoa-fisica');
    
    tipoPessoaSelect.addEventListener('change', function() {
        const tipo = this.value;
        
        if (tipo === 'fisica') {
            documentoLabel.innerHTML = '<i class="bi bi-card-text me-2"></i>CPF';
            nomeLabel.innerHTML = '<i class="bi bi-person me-2"></i>Nome Completo';
            nomeFantasiaContainer.classList.add('d-none');
            documentoInput.placeholder = '000.000.000-00';
            documentoInput.maxLength = 14;
            
            // Mostrar campos específicos para pessoa física
            pessoaFisicaElements.forEach(el => {
                el.classList.remove('d-none');
                const inputs = el.querySelectorAll('input');
                inputs.forEach(input => input.setAttribute('required', ''));
            });
        } else if (tipo === 'juridica') {
            documentoLabel.innerHTML = '<i class="bi bi-building me-2"></i>CNPJ';
            nomeLabel.innerHTML = '<i class="bi bi-building me-2"></i>Razão Social';
            nomeFantasiaContainer.classList.remove('d-none');
            documentoInput.placeholder = '00.000.000/0000-00';
            documentoInput.maxLength = 18;
            
            // Esconder campos específicos para pessoa física
            pessoaFisicaElements.forEach(el => {
                el.classList.add('d-none');
                const inputs = el.querySelectorAll('input');
                inputs.forEach(input => input.removeAttribute('required'));
            });
        }
        
        // Limpar o campo de documento
        documentoInput.value = '';
    });
    
    // Evento para formatar e validar documento
    documentoInput.addEventListener('input', function() {
        const tipo = tipoPessoaSelect.value;
        let valor = this.value.replace(/\D/g, '');
        
        if (tipo === 'fisica') {
            if (valor.length <= 11) {
                this.value = formatCPF(valor);
            }
        } else if (tipo === 'juridica') {
            if (valor.length <= 14) {
                this.value = formatCNPJ(valor);
            }
        }
    });
});

// Submissão do formulário
document.getElementById('registerForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    if (!this.checkValidity()) {
        e.stopPropagation();
        this.classList.add('was-validated');
        return;
    }
    
    const submitBtn = this.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    
    // Loading state
    submitBtn.classList.add('btn-loading');
    submitBtn.disabled = true;
    
    try {
        // Simular chamada API
        await new Promise(resolve => setTimeout(resolve, 3000));
        
        EvoDrones.showToast('Conta criada com sucesso! Redirecionando...', 'success');
        
        setTimeout(() => {
            window.location.href = '../login/';
        }, 2000);
        
    } catch (error) {
        EvoDrones.showToast('Erro ao criar conta. Tente novamente.', 'danger');
    } finally {
        submitBtn.classList.remove('btn-loading');
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }
});
 
 // Evento para consultar CNPJ quando o campo perder o foco
 document.getElementById('documento').addEventListener('blur', async function() {
     const tipoPessoaSelect = document.getElementById('tipoPessoa');
     const tipo = tipoPessoaSelect.value;
     
     if (tipo === 'juridica') {
         const cnpj = this.value.replace(/\D/g, '');
         
         if (cnpj.length === 14 && isValidCNPJ(cnpj)) {
             // Mostrar indicador de carregamento
             const loadingIndicator = document.createElement('div');
             loadingIndicator.id = 'cnpjLoading';
             loadingIndicator.className = 'text-center mt-2';
             loadingIndicator.innerHTML = '<div class="spinner-border spinner-border-sm text-primary" role="status"></div> <span>Consultando CNPJ...</span>';
             
             this.parentNode.after(loadingIndicator);
             
             // Consultar CNPJ
             const dadosCNPJ = await consultarCNPJ(cnpj);
             
             // Remover indicador de carregamento
             document.getElementById('cnpjLoading')?.remove();
             
             if (dadosCNPJ) {
                 // Preencher os campos com os dados retornados
                 document.getElementById('nome').value = dadosCNPJ.razao_social || '';
                 document.getElementById('nomeFantasia').value = dadosCNPJ.nome_fantasia || '';
                 document.getElementById('email').value = dadosCNPJ.email || '';
                 document.getElementById('telefone').value = dadosCNPJ.ddd_telefone_1 || '';
             }
         }
     }
 });