function scrollSpy() {
    return {
        activeSection: 'hero',
        sections: ['features', 'architecture', 'providers', 'installation'],
        init() {
            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting) {
                        this.activeSection = entry.target.id;
                    }
                });
            }, { threshold: 0.3 });
            
            this.sections.forEach(id => {
                const el = document.getElementById(id);
                if (el) observer.observe(el);
            });
        }
    };
}
