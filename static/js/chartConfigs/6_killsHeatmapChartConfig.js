// static/js/chartConfigs/killsHeatmapChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kills Heatmap Chart
 */
const killsHeatmapChartConfig = {
    type: 'matrix', // Using matrix chart type from chartjs-chart-matrix
    options: getCommonOptions('Kills Heatmap', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    title: function () { return ''; }, // Remove title
                    label: function (context) {
                        const xLabel = context.raw.x || 'Unknown';
                        const yLabel = context.raw.y || 'Unknown';
                        const value = context.raw.v || 0;
                        return `Day: ${yLabel}, Hour: ${xLabel}, Kills: ${value}`;
                    },
                },
            },
            datalabels: {
                display: false, // Typically not needed for heatmaps
            },
        },
        scales: {
            x: {
                type: 'category',
                labels: [
                    '0', '1', '2', '3', '4', '5', '6',
                    '7', '8', '9', '10', '11', '12', '13',
                    '14', '15', '16', '17', '18', '19', '20',
                    '21', '22', '23'
                ],
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
                title: {
                    display: true,
                    text: 'Hour of Day',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                type: 'category',
                labels: [
                    'Sunday', 'Monday', 'Tuesday', 'Wednesday',
                    'Thursday', 'Friday', 'Saturday'
                ],
                ticks: {
                    color: '#ffffff',
                },
                grid: { display: false },
                title: {
                    display: true,
                    text: 'Day of Week',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
        },
        elements: {
            rectangle: {
                borderWidth: 0, // Remove borders
                // borderColor: '#ffffff', // Optional: Remove or comment out
                // borderSkipped: 'bottom', // Optional: Remove or comment out
            },
        },
        responsive: true,
        maintainAspectRatio: false,
    }),
    processData: function (data) {
        const chartName = 'Kills Heatmap';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty data to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        // console.log('Incoming data for Kills Heatmap:', data); // Debugging log

        // Ensure data has exactly 7 days
        if (data.length !== 7) {
            console.warn(`Kills Heatmap expects data for 7 days, received ${data.length}.`);
        }

        // Prepare matrix data
        const matrixData = [];
        const yLabels = [
            'Sunday', 'Monday', 'Tuesday', 'Wednesday',
            'Thursday', 'Friday', 'Saturday'
        ];
        const xLabels = [
            '0', '1', '2', '3', '4', '5', '6',
            '7', '8', '9', '10', '11', '12', '13',
            '14', '15', '16', '17', '18', '19', '20',
            '21', '22', '23'
        ];

        let maxKills = 0; // To determine scaling for backgroundColor

        for (let day = 0; day < yLabels.length; day++) {
            for (let hour = 0; hour < xLabels.length; hour++) {
                const kills = (data[day] && typeof data[day][hour] === 'number') ? data[day][hour] : 0;
                matrixData.push({
                    x: xLabels[hour],
                    y: yLabels[day],
                    v: kills,
                });
                if (kills > maxKills) {
                    maxKills = kills;
                }
            }
        }

        // Check if there are at least 3 days with data
        let daysWithData = 0;
        for (let day = 0; day < yLabels.length; day++) {
            const totalKillsPerDay = (data[day] && Array.isArray(data[day])) ? data[day].reduce((a, b) => a + (b || 0), 0) : 0;
            if (totalKillsPerDay > 0) {
                daysWithData++;
            }
        }

        if (daysWithData < 3) {
            console.warn(`Not enough data points (${daysWithData} days) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        // Generate background colors array
        const backgroundColors = matrixData.map(dataPoint => {
            const kills = dataPoint.v;
            if (kills > 0) {
                // Calculate opacity based on kill count
                const alpha = maxKills > 0 ? Math.min(kills / maxKills, 1) : 0;
                // Return shades of red with opacity
                return `rgba(255, 0, 0, ${0.3 + 0.7 * alpha})`;
            } else {
                return 'rgba(0, 0, 0, 0)'; // Fully transparent
            }
        });

        // Debugging logs
        // console.log('Matrix Data:', matrixData);
        // console.log('Max Kills:', maxKills);
        // console.log('Background Colors:', backgroundColors);
        // console.log('Border Colors:', borderColors);

        const dataset = {
            label: 'Kills Heatmap',
            data: matrixData,
            backgroundColor: backgroundColors, // Array of colors
        };

        // console.log('Processed dataset for Kills Heatmap:', dataset); // Debugging log

        return { labels: [], datasets: [dataset] };
    },
};

export default killsHeatmapChartConfig;
