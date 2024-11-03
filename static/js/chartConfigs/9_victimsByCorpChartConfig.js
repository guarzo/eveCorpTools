// static/js/chartConfigs/9_victimsByCorpChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Victims by Corporation Chart
 */
const victimsByCorporationChartConfig = {
    type: 'bar',
    options: getCommonOptions('Victims by Corporation', {
        scales: {
            x: {
                title: {
                    display: true,
                    text: 'Corporation',
                },
                ticks: {
                    color: '#ffffff',
                    autoSkip: false, // Show all labels
                    maxRotation: 90,
                    minRotation: 45,
                },
                grid: { display: false },
            },
            y: {
                title: {
                    display: true,
                    text: 'Number of Victims',
                },
                ticks: {
                    color: '#ffffff',
                    beginAtZero: true,
                },
                grid: { display: true, color: '#444444' },
            },
        },
        plugins: {
            legend: {
                display: false, // Single dataset; legend not needed
            },
            tooltip: {
                callbacks: {
                    label: function(context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y !== null ? context.parsed.y.toLocaleString() : '0';
                        return `${label}: ${value}`;
                    },
                },
            },
        },
        responsive: true,
        maintainAspectRatio: false,
    }),
    processData: function(data) {
        const chartName = 'Victims by Corporation';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        // console.log('Incoming data for Victims by Corporation:', data); // Debugging log
        //
        // // Inspect each data item
        // data.forEach((item, index) => {
        //     console.log(`Item ${index}:`, item);
        // });

        // Map the correct fields
        const labels = data.map(item => item.name || 'Unknown');
        const victims = data.map(item => item.kill_count || 0);

        // Check for 'Unknown' labels
        const allUnknown = labels.every(label => label === 'Unknown');
        if (allUnknown) {
            console.warn(`All labels for ${chartName} are 'Unknown'. Check data source.`);
        }

        // Number of bars
        const count = labels.length;

        // Generate distinct, random colors
        const backgroundColors = generateDistinctColors(count);
        const borderColors = backgroundColors.map(color => {
            // Convert backgroundColor to a fully opaque version for borders
            return color.replace(/rgba\((\d+),\s*(\d+),\s*(\d+),\s*[\d.]+\)/, 'rgba($1, $2, $3, 1)');
        });

        return {
            labels: labels,
            datasets: [{
                label: 'Number of Victims',
                data: victims,
                backgroundColor: backgroundColors,
                borderColor: borderColors,
                borderWidth: 1,
            }]
        };

        /**
         * Generates an array of distinct, random RGBA colors.
         * @param {number} count - Number of colors to generate.
         * @returns {string[]} Array of RGBA color strings.
         */
        function generateDistinctColors(count) {
            const colors = [];
            const saturation = 70; // Percentage
            const lightness = 50; // Percentage

            for (let i = 0; i < count; i++) {
                const hue = Math.floor(Math.random() * 360); // Random hue between 0 and 359
                const alpha = 0.6; // Opacity for background
                colors.push(`rgba(${hslToRgb(hue, saturation, lightness)}, ${alpha})`);
            }

            // Ensure colors are unique
            return ensureUniqueColors(colors);
        }

        /**
         * Converts HSL to RGB.
         * @param {number} h - Hue (0-360).
         * @param {number} s - Saturation (0-100).
         * @param {number} l - Lightness (0-100).
         * @returns {string} RGB string formatted as 'r, g, b'.
         */
        function hslToRgb(h, s, l) {
            s /= 100;
            l /= 100;

            const c = (1 - Math.abs(2 * l - 1)) * s;
            const hh = h / 60;
            const x = c * (1 - Math.abs((hh % 2) - 1));

            let r = 0, g = 0, b = 0;

            if (0 <= hh && hh < 1) {
                r = c; g = x; b = 0;
            } else if (1 <= hh && hh < 2) {
                r = x; g = c; b = 0;
            } else if (2 <= hh && hh < 3) {
                r = 0; g = c; b = x;
            } else if (3 <= hh && hh < 4) {
                r = 0; g = x; b = c;
            } else if (4 <= hh && hh < 5) {
                r = x; g = 0; b = c;
            } else if (5 <= hh && hh < 6) {
                r = c; g = 0; b = x;
            }

            const m = l - c / 2;
            r = Math.round((r + m) * 255);
            g = Math.round((g + m) * 255);
            b = Math.round((b + m) * 255);

            return `${r}, ${g}, ${b}`;
        }

        /**
         * Ensures all colors in the array are unique. If duplicates are found, regenerate until unique.
         * @param {string[]} colors - Array of RGBA color strings.
         * @returns {string[]} Array of unique RGBA color strings.
         */
        function ensureUniqueColors(colors) {
            const uniqueColors = [];
            const colorSet = new Set();

            for (let color of colors) {
                // Simple uniqueness check based on color string
                if (!colorSet.has(color)) {
                    colorSet.add(color);
                    uniqueColors.push(color);
                } else {
                    // Regenerate a unique color if duplicate is found
                    let newColor;
                    do {
                        const hue = Math.floor(Math.random() * 360);
                        const saturation = 70;
                        const lightness = 50;
                        const alpha = 0.6;
                        newColor = `rgba(${hslToRgb(hue, saturation, lightness)}, ${alpha})`;
                    } while (colorSet.has(newColor));
                    colorSet.add(newColor);
                    uniqueColors.push(newColor);
                }
            }

            return uniqueColors;
        }
    },
};

export default victimsByCorporationChartConfig;
