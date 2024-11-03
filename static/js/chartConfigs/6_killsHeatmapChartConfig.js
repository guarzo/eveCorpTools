// static/js/chartConfigs/6_killsHeatmapChartConfig.js

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
                        const xLabel = context.raw.x; // Access x value directly
                        const yLabel = context.raw.y; // Access y value directly
                        const value = context.raw.v;
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
                borderWidth: 1,
                borderColor: '#ffffff',
                borderSkipped: 'bottom',
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

        console.log('Incoming data for Kills Heatmap:', data); // Debugging log

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

        for (let day = 0; day < 7; day++) {
            for (let hour = 0; hour < 24; hour++) {
                const kills = (data[day] && data[day][hour] !== undefined) ? data[day][hour] : 0;
                matrixData.push({
                    x: hour.toString(),
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
        for (let day = 0; day < 7; day++) {
            const totalKillsPerDay = (data[day] && Array.isArray(data[day])) ? data[day].reduce((a, b) => a + (b || 0), 0) : 0;
            if (totalKillsPerDay > 0) {
                daysWithData++;
            }
        }

        if (daysWithData < 3) {
            console.warn(`Not enough data points (${daysWithData} days) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        const dataset = {
            label: 'Kills Heatmap',
            data: matrixData,
            // Removed xLabels and yLabels from dataset as they are not needed
            backgroundColor: function (context) {
                const value = context.dataset.data[context.dataIndex].v;
                if (value > 0) {
                    // Define a red color with varying opacity based on kill count
                    const alpha = maxKills > 0 ? Math.min(value / maxKills, 1) : 0;
                    return `rgba(255, 0, 0, ${alpha})`; // Red with varying opacity
                } else {
                    // Assign a light gray color for no kills
                    return 'rgba(211, 211, 211, 0.1)'; // Light gray with low opacity
                }
            },
        };

        return { labels: [], datasets: [dataset] };
    },
};

export default killsHeatmapChartConfig;
