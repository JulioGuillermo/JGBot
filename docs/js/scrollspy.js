function scrollSpy() {
    return {
        activeSection: 'hero',
        sections: ['features', 'architecture', 'providers', 'installation'],
        init() {
            // Use scroll event to determine which section is most visible
            // This handles small sections that don't trigger IntersectionObserver properly
            window.addEventListener('scroll', () => this.updateActiveSection());
            // Initial check
            this.updateActiveSection();
        },
        updateActiveSection() {
            const viewportHeight = window.innerHeight;
            const scrollTop = window.scrollY;
            
            let maxVisibility = 0;
            let mostVisibleSection = this.activeSection;
            
            this.sections.forEach(id => {
                const el = document.getElementById(id);
                if (!el) return;
                
                const rect = el.getBoundingClientRect();
                const sectionTop = rect.top;
                const sectionHeight = rect.height;
                
                // Calculate how much of the section is visible in the viewport
                const visibleTop = Math.max(0, sectionTop);
                const visibleBottom = Math.min(viewportHeight, sectionTop + sectionHeight);
                const visibleHeight = Math.max(0, visibleBottom - visibleTop);
                const visibility = visibleHeight / sectionHeight;
                
                // Calculate center proximity - sections closer to viewport center get priority
                const sectionCenter = sectionTop + sectionHeight / 2;
                const viewportCenter = viewportHeight / 2;
                const distanceFromCenter = Math.abs(sectionCenter - viewportCenter);
                const centerProximity = 1 - (distanceFromCenter / (viewportHeight / 2));
                
                // Combined score: visibility + center proximity
                const score = visibility * 0.6 + centerProximity * 0.4;
                
                // Only switch if the new section has significantly higher score
                // or if current section is no longer visible
                const currentEl = document.getElementById(this.activeSection);
                const isCurrentVisible = currentEl && 
                    currentEl.getBoundingClientRect().top < viewportHeight && 
                    currentEl.getBoundingClientRect().bottom > 0;
                
                if (score > maxVisibility && (!isCurrentVisible || score > maxVisibility + 0.1)) {
                    maxVisibility = score;
                    mostVisibleSection = id;
                }
            });
            
            if (mostVisibleSection !== this.activeSection && maxVisibility > 0.1) {
                this.activeSection = mostVisibleSection;
            }
        }
    };
}
