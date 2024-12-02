// Theme switcher
const themeSwitcher = {
  // Config
  _scheme : 'auto',
  menuTarget : "details[role='list']",
  buttonsTarget : 'a[data-theme-switcher]',
  buttonAttribute : 'data-theme-switcher',
  rootAttribute : 'data-theme',
  localStorageKey : 'theme',

  // Init theme
  initTheme() {
    this._scheme = this.schemeFromLocalStorage || this.preferredColorScheme;
    this.scheme = this._scheme;
  },

  get schemeFromLocalStorage() {
    return localStorage.getItem(this.localStorageKey);
  },

  get preferredColorScheme() {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark'
                                                                     : 'light';
  },

  // Set theme
  set scheme(scheme) {
    if (scheme === 'auto') {
      this._scheme = this.preferredColorScheme;
    } else {
      this._scheme = scheme;
    }
    this.applyScheme();
    this.schemeToLocalStorage();
  },

  // Get theme
  get scheme() { return this._scheme; },

  // Apply theme
  applyScheme() {
    document.querySelector('html').setAttribute(this.rootAttribute,
                                                this._scheme);
  },

  // Store theme
  schemeToLocalStorage() {
    localStorage.setItem(this.localStorageKey, this._scheme);
  },

  // Click handler
  handleClick(event) {
    if (event.target.matches(this.buttonsTarget)) {
      event.preventDefault();
      this.scheme = event.target.getAttribute(this.buttonAttribute);
    }
  },
};

// Init theme switcher
themeSwitcher.initTheme();

// Theme switcher event listener
document.addEventListener('click',
                          (event) => { themeSwitcher.handleClick(event); });
