async function copyToClipboard(text) {
    const button = document.querySelector('.install-link-container button');

    try {
        await navigator.clipboard.writeText(text);
        button.classList.add('copied-tooltip');
        setTimeout(() => {
            button.classList.remove('copied-tooltip');
        }, 500);
        console.log('Text copied to clipboard');
    } catch (err) {
        console.error('Failed to copy text: ', err);
    }
}

// Disk Usage Visualization Interactivity
document.addEventListener('DOMContentLoaded', () => {
    const barSegments = document.querySelectorAll('.bar-segment');
    const diskUsageContainer = document.querySelector('.disk-usage-container');

    if (barSegments.length > 0 && diskUsageContainer) {
        barSegments.forEach(segment => {
            let tooltip = document.createElement('div');
            tooltip.classList.add('disk-usage-tooltip');
            diskUsageContainer.appendChild(tooltip);

            segment.addEventListener('mousemove', (e) => {
                tooltip.style.display = 'block';
                const segmentType = segment.classList[1];
                let size = 'N/A';

                switch (segmentType) {
                    case 'apps':
                        size = '41.67 GB';
                        tooltip.innerHTML = `
                            <strong>${segmentType.charAt(0).toUpperCase() + segmentType.slice(1)}</strong><br>
                            ${size}
                        `;
                        break;
                    case 'documents':
                        size = '36.76 GB';
                        tooltip.innerHTML = `
                            <strong>${segmentType.charAt(0).toUpperCase() + segmentType.slice(1)}</strong><br>
                            ${size}
                        `;
                        break;
                    case 'developer':
                        size = '7.35 GB';
                        tooltip.innerHTML = `
                            <strong>${segmentType.charAt(0).toUpperCase() + segmentType.slice(1)}</strong><br>
                            ${size}
                        `;
                        break;
                    case 'macos':
                        size = '25.7 GB';
                        tooltip.innerHTML = `
                            <strong>${segmentType.charAt(0).toUpperCase() + segmentType.slice(1)}</strong><br>
                            ${size}
                        `;
                        break;
                    case 'system-data':
                        size = '122.55 GB';
                        tooltip.innerHTML = `
                            <strong>${segmentType.charAt(0).toUpperCase() + segmentType.slice(1)}</strong><br>
                            ${size}
                        `;
                        tooltip.innerHTML += `<br>What the fuck?!!`; // Append the additional text here
                        break;
                }

                // Position the tooltip relative to the segment
                const containerRect = diskUsageContainer.getBoundingClientRect();
                const segmentRect = segment.getBoundingClientRect();
                tooltip.style.display = 'block'; // Ensure tooltip is visible to get its dimensions
                const tooltipRect = tooltip.getBoundingClientRect();

                // Calculate x position to center the tooltip above the segment
                let x = (segmentRect.left + segmentRect.width / 2) - containerRect.left - tooltipRect.width / 15;
                // Calculate y position above the segment with a small offset
                let y = segmentRect.top - containerRect.top - tooltipRect.height - 5;

                // Ensure tooltip stays within container bounds
                if (x + tooltipRect.width > containerRect.width) {
                    x = containerRect.width - tooltipRect.width - 5; // 5px padding from right edge
                }
                if (x < 0) {
                    x = 5; // 5px padding from left edge
                }
                // Re-evaluate y-axis boundary check
                if (y < 0) {
                    // If not enough space above, position below the segment
                    y = segmentRect.bottom - containerRect.top + 5;
                    // If still not enough space (e.g., very small container), position at top with padding
                    if (y + tooltipRect.height > containerRect.height) {
                        y = 5;
                    }
                }

                tooltip.style.left = `${x}px`;
                tooltip.style.top = `${y}px`;
            });

            segment.addEventListener('mouseleave', () => {
                tooltip.style.display = 'none';
            });
        });
    }
});