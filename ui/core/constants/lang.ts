// Returns a 2-letter browser language code, or '' if not available.
export function getBrowserLanguageCode(): string {
	return (navigator.language || '').substring(0, 2);
}

export function getLanguageCode(): string {
	return cachedLanguageCode_;
}

export function setLanguageCode(newLang: string) {
	// Keep language code for sim settings; 'en' is stored as empty string for compatibility
	cachedLanguageCode_ = newLang == 'en' ? '' : newLang;
}

let cachedLanguageCode_ = '';
