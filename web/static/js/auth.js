// Auth Page JavaScript - Film Matching App
// Includes tab switching, form validation, animations, and API integration placeholders

class AuthManager {
    constructor() {
        this.currentTab = 'login';
        this.isLoading = false;
        this.validators = {
            username: (value) => value.length >= 3,
            email: (value) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value),
            password: (value) => value.length >= 8,
            confirmPassword: (value, password) => value === password
        };

        this.init();
        this.initInitialAnimation();
    }

    init() {
        this.initTabSwitching();
        this.initFormValidation();
        this.initPasswordToggles();
        this.initSubmitHandlers();
        this.initAnimations();
    }

    initInitialAnimation() {
        // Добавляем плавную анимацию появления для изначально видимой формы логина
        const loginForm = document.getElementById('login-form');
        loginForm.style.opacity = '0';
        loginForm.style.transform = 'translateY(20px)';

        // Небольшая задержка для плавности
        setTimeout(() => {
            loginForm.style.transition = 'opacity 0.5s ease-out, transform 0.5s ease-out';
            loginForm.style.opacity = '1';
            loginForm.style.transform = 'translateY(0)';
        }, 100);
    }

    // Tab Switching Logic
    initTabSwitching() {
        const tabButtons = document.querySelectorAll('.tab-button');
        const tabIndicator = document.querySelector('.tab-indicator');
        const loginForm = document.getElementById('login-form');
        const registerForm = document.getElementById('register-form');

        tabButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                const tabType = e.target.dataset.tab;
                this.switchTab(tabType, tabButtons, tabIndicator, loginForm, registerForm);
            });
        });
    }

    switchTab(tabType, tabButtons, tabIndicator, loginForm, registerForm) {
        if (this.currentTab === tabType || this.isLoading) return;

        // Update tab buttons
        tabButtons.forEach(btn => btn.classList.remove('active'));
        document.querySelector(`[data-tab="${tabType}"]`).classList.add('active');

        // Animate tab indicator
        if (tabType === 'register') {
            tabIndicator.classList.add('register');
        } else {
            tabIndicator.classList.remove('register');
        }

        // Animate form switching
        const currentForm = this.currentTab === 'login' ? loginForm : registerForm;
        const nextForm = tabType === 'login' ? loginForm : registerForm;

        this.animateFormSwitch(currentForm, nextForm);
        this.currentTab = tabType;

        // Clear any validation errors
        this.clearFormErrors();
    }

    animateFormSwitch(currentForm, nextForm) {
        // Очищаем любые inline стили перед анимацией
        currentForm.style.transition = '';
        currentForm.style.opacity = '';
        currentForm.style.transform = '';
        nextForm.style.transition = '';
        nextForm.style.opacity = '';
        nextForm.style.transform = '';

        // Slide out current form
        currentForm.classList.add('slide-out');

        setTimeout(() => {
            currentForm.classList.add('hidden');
            currentForm.classList.remove('slide-out');

            // Slide in next form
            nextForm.classList.remove('hidden');
            nextForm.classList.add('slide-in');

            setTimeout(() => {
                nextForm.classList.remove('slide-in');
            }, 300);
        }, 150);
    }

    // Form Validation
    initFormValidation() {
        const inputs = document.querySelectorAll('input');

        inputs.forEach(input => {
            input.addEventListener('blur', () => this.validateField(input));
            input.addEventListener('input', () => {
                this.clearFieldError(input);

                // Если изменился пароль, проверяем также confirmPassword в той же форме
                if (input.name === 'password') {
                    const form = input.closest('form');
                    const confirmPasswordField = form.querySelector('input[name="confirmPassword"]');
                    if (confirmPasswordField && confirmPasswordField.value) {
                        // Небольшая задержка для лучшего UX
                        setTimeout(() => this.validateField(confirmPasswordField), 100);
                    }
                }
            });
        });
    }

    validateField(input) {
        const inputGroup = input.closest('.input-group');
        const fieldName = input.name;
        const value = input.value.trim();

        let isValid = true;
        let errorMessage = '';

        switch (fieldName) {
            case 'name':
                isValid = this.validators.username(value); // используем ту же валидацию для имени
                errorMessage = 'Имя должно содержать минимум 3 символа';
                break;

            case 'email':
                isValid = this.validators.email(value);
                errorMessage = 'Введите корректный email адрес';
                break;

            case 'password':
                isValid = this.validators.password(value);
                errorMessage = 'Пароль должен содержать минимум 8 символов';
                break;

            case 'confirmPassword':
                // Ищем поле password в той же форме, где находится confirmPassword
                const form = input.closest('form');
                const passwordField = form.querySelector('input[name="password"]');
                isValid = this.validators.confirmPassword(value, passwordField ? passwordField.value : '');
                errorMessage = 'Пароли не совпадают';
                break;
        }

        if (value && !isValid) {
            this.showFieldError(inputGroup, errorMessage);
        } else if (value && isValid) {
            this.showFieldSuccess(inputGroup);
        }

        return isValid;
    }

    showFieldError(inputGroup, message) {
        inputGroup.classList.remove('success');
        inputGroup.classList.add('error');

        let errorElement = inputGroup.querySelector('.error-message');
        if (!errorElement) {
            errorElement = document.createElement('div');
            errorElement.className = 'error-message';
            inputGroup.appendChild(errorElement);
        }
        errorElement.textContent = message;
    }

    showFieldSuccess(inputGroup) {
        inputGroup.classList.remove('error');
        inputGroup.classList.add('success');

        const errorElement = inputGroup.querySelector('.error-message');
        if (errorElement) {
            errorElement.remove();
        }
    }

    clearFieldError(input) {
        const inputGroup = input.closest('.input-group');
        inputGroup.classList.remove('error', 'success');

        const errorElement = inputGroup.querySelector('.error-message');
        if (errorElement) {
            errorElement.remove();
        }
    }

    clearFormErrors() {
        const errorInputs = document.querySelectorAll('.input-group.error, .input-group.success');
        errorInputs.forEach(inputGroup => {
            inputGroup.classList.remove('error', 'success');
            const errorElement = inputGroup.querySelector('.error-message');
            if (errorElement) {
                errorElement.remove();
            }
        });
    }

    // Password Toggle Functionality
    initPasswordToggles() {
        const passwordToggles = document.querySelectorAll('.password-toggle');

        passwordToggles.forEach(toggle => {
            toggle.addEventListener('click', (e) => {
                const input = e.target.closest('.input-group').querySelector('input');
                const isPassword = input.type === 'password';

                input.type = isPassword ? 'text' : 'password';

                // Update icon (you can customize this based on your icon library)
                const svg = toggle.querySelector('svg');
                if (isPassword) {
                    // Eye-off icon (password visible)
                    svg.innerHTML = `
                        <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                        <line x1="1" y1="1" x2="23" y2="23"/>
                    `;
                } else {
                    // Eye icon (password hidden)
                    svg.innerHTML = `
                        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                        <circle cx="12" cy="12" r="3"/>
                    `;
                }
            });
        });
    }

    // Form Submit Handlers
    initSubmitHandlers() {
        const loginForm = document.getElementById('loginForm');
        const registerForm = document.getElementById('registerForm');

        loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        registerForm.addEventListener('submit', (e) => this.handleRegister(e));
    }

    async handleLogin(e) {
        e.preventDefault();

        if (this.isLoading) return;

        const formData = new FormData(e.target);
        const data = {
            email: formData.get('email'),
            password: formData.get('password'),
            remember: formData.get('remember') === 'on'
        };

        // Validate form
        if (!this.validateLoginForm(data)) {
            return;
        }

        this.setLoading(true, 'loginButton');

        try {
            // TODO: Replace with actual API call
            const result = await this.loginAPI(data);

            if (result.success) {
                this.showSuccess('Вход выполнен успешно!');
                // Переадресация после успешного логина
                setTimeout(() => {
                    window.location.href = 'home.html';
                }, 1500);
            } else {
                this.showError(result.message || 'Ошибка входа');
            }
        } catch (error) {
            console.error('Login error:', error);
            this.showError('Произошла ошибка. Попробуйте позже.');
        } finally {
            this.setLoading(false, 'loginButton');
        }
    }

    async handleRegister(e) {
        e.preventDefault();

        if (this.isLoading) return;

        const formData = new FormData(e.target);
        const data = {
            name: formData.get('name'),
            email: formData.get('email'),
            password: formData.get('password'),
            confirmPassword: formData.get('confirmPassword'),
            terms: formData.get('terms') === 'on'
        };

        // Validate form
        if (!this.validateRegisterForm(data)) {
            return;
        }

        this.setLoading(true, 'registerButton');

        try {
            // TODO: Replace with actual API call
            const result = await this.registerAPI(data);

            if (result.success) {
                this.showSuccess('Аккаунт создан успешно! Теперь войдите в систему.');
                // Переключаем на таб логина после успешной регистрации
                setTimeout(() => {
                    const tabButtons = document.querySelectorAll('.tab-button');
                    const tabIndicator = document.querySelector('.tab-indicator');
                    const loginForm = document.getElementById('login-form');
                    const registerForm = document.getElementById('register-form');
                    this.switchTab('login', tabButtons, tabIndicator, loginForm, registerForm);
                }, 2000);
            } else {
                this.showError(result.message || 'Ошибка регистрации');
            }
        } catch (error) {
            console.error('Register error:', error);
            this.showError('Произошла ошибка. Попробуйте позже.');
        } finally {
            this.setLoading(false, 'registerButton');
        }
    }

    validateLoginForm(data) {
        let isValid = true;

        if (!data.email) {
            this.showError('Введите email');
            isValid = false;
        } else if (!this.validators.email(data.email)) {
            this.showError('Введите корректный email адрес');
            isValid = false;
        }

        if (!data.password) {
            this.showError('Введите пароль');
            isValid = false;
        }

        return isValid;
    }

    validateRegisterForm(data) {
        let isValid = true;

        if (!this.validators.username(data.name)) {
            this.showError('Имя должно содержать минимум 3 символа');
            isValid = false;
        }

        if (!this.validators.email(data.email)) {
            this.showError('Введите корректный email адрес');
            isValid = false;
        }

        if (!this.validators.password(data.password)) {
            this.showError('Пароль должен содержать минимум 8 символов');
            isValid = false;
        }

        if (!this.validators.confirmPassword(data.confirmPassword, data.password)) {
            this.showError('Пароли не совпадают');
            isValid = false;
        }

        if (!data.terms) {
            this.showError('Необходимо согласиться с условиями использования');
            isValid = false;
        }

        return isValid;
    }

    // Loading State Management
    setLoading(loading, buttonId) {
        this.isLoading = loading;
        const button = document.getElementById(buttonId);

        if (loading) {
            button.classList.add('loading');
            button.disabled = true;
        } else {
            button.classList.remove('loading');
            button.disabled = false;
        }
    }

    // Notification System
    showSuccess(message) {
        this.showNotification(message, 'success');
    }

    showError(message) {
        this.showNotification(message, 'error');
    }

    showNotification(message, type) {
        // Remove existing notifications
        const existing = document.querySelector('.notification');
        if (existing) {
            existing.remove();
        }

        // Create notification element
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.innerHTML = `
            <div class="notification-content">
                <span class="notification-message">${message}</span>
                <button class="notification-close" onclick="this.parentElement.parentElement.remove()">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>
            </div>
        `;

        // Add styles
        notification.style.cssText = `
            position: fixed;
            top: 2rem;
            right: 2rem;
            background: ${type === 'success' ? '#04d361' : '#e87c03'};
            color: white;
            padding: 1rem 1.5rem;
            border-radius: 8px;
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
            z-index: 1000;
            animation: slideInFromRight 0.3s ease-out;
            max-width: 400px;
        `;

        document.body.appendChild(notification);

        // Auto remove after 5 seconds
        setTimeout(() => {
            if (notification.parentElement) {
                notification.style.animation = 'slideOutToRight 0.3s ease-in forwards';
                setTimeout(() => notification.remove(), 300);
            }
        }, 5000);
    }

    // API Methods - Real API calls to Go backend
    async loginAPI(data) {
        try {
            const requestData = {
                email: data.email,
                password: data.password
            };

            console.log('Отправляем данные логина:', requestData);
            console.log('JSON строка:', JSON.stringify(requestData));

            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Добавляем для сохранения cookie
                body: JSON.stringify(requestData)
            });

            console.log('Статус ответа:', response.status, response.statusText);
            console.log('Заголовки ответа:', [...response.headers.entries()]);

            if (response.ok) {
                // Успешный логин - cookie установлен автоматически
                console.log('Логин успешен');
                return {success: true, message: 'Вход выполнен успешно!'};
            } else {
                // Ошибка - получаем текст ошибки
                const errorText = await response.text();
                console.log('Ошибка логина:', errorText);
                return {success: false, message: errorText || 'Ошибка входа'};
            }
        } catch (error) {
            console.error('Login API error:', error);
            return {success: false, message: 'Ошибка сети. Проверьте соединение.'};
        }
    }

    async registerAPI(data) {
        try {
            const requestData = {
                name: data.name,
                email: data.email,
                password: data.password
            };

            console.log('Отправляем данные регистрации:', requestData);
            console.log('JSON строка:', JSON.stringify(requestData));

            const response = await fetch('http://localhost:8080/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include', // Добавляем для консистентности
                body: JSON.stringify(requestData)
            });

            console.log('Статус ответа:', response.status, response.statusText);
            console.log('Заголовки ответа:', [...response.headers.entries()]);

            if (response.ok) {
                // Успешная регистрация
                console.log('Регистрация успешна');
                return {success: true, message: 'Аккаунт создан успешно!'};
            } else {
                // Ошибка - получаем текст ошибки
                const errorText = await response.text();
                console.log('Ошибка регистрации:', errorText);
                return {success: false, message: errorText || 'Ошибка регистрации'};
            }
        } catch (error) {
            console.error('Register API error:', error);
            return {success: false, message: 'Ошибка сети. Проверьте соединение.'};
        }
    }

    // Animation Helpers
    initAnimations() {
        // Add CSS for notifications if not already present
        if (!document.querySelector('#notification-styles')) {
            const style = document.createElement('style');
            style.id = 'notification-styles';
            style.textContent = `
                @keyframes slideInFromRight {
                    from {
                        opacity: 0;
                        transform: translateX(100%);
                    }
                    to {
                        opacity: 1;
                        transform: translateX(0);
                    }
                }

                @keyframes slideOutToRight {
                    from {
                        opacity: 1;
                        transform: translateX(0);
                    }
                    to {
                        opacity: 0;
                        transform: translateX(100%);
                    }
                }

                .notification-content {
                    display: flex;
                    align-items: center;
                    gap: 1rem;
                }

                .notification-close {
                    background: none;
                    border: none;
                    color: inherit;
                    cursor: pointer;
                    padding: 0.25rem;
                    border-radius: 4px;
                    transition: background 0.2s ease;
                }

                .notification-close:hover {
                    background: rgba(255, 255, 255, 0.2);
                }
            `;
            document.head.appendChild(style);
        }
    }
}

// Initialize the auth manager when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new AuthManager();
});

// Export for potential external use
window.AuthManager = AuthManager;
