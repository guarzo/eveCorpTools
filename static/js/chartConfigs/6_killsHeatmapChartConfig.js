// static/js/chartConfigs/killsHeatmapChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kills Heatmap Chart
 */
const killsHeatmapChartConfig = {
    id: 'killsHeatmapChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdKillHeatmapData', canvasId: 'killsHeatmapChart_mtd' },
        ytd: { dataVar: 'ytdKillHeatmapData', canvasId: 'killsHeatmapChart_ytd' },
        lastMonth: { dataVar: 'lastMKillHeatmapData', canvasId: 'killsHeatmapChart_lastM' },
    },
    type: 'matrix', // Using matrix chart type from chartjs-chart-matrix
    dataType: 'array', // Specify that this chart expects array data
    options: getCommonOptions('Kills Heatmap', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    title: function () { return ''; }, // Remove title
                    label: function (context) {
                        const xLabel = context.dataset.xLabels[context.dataIndex];
                        const yLabel = context.dataset.yLabels[context.dataIndex];
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
                const kills = data[day][hour] || 0;
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

        // Check if there are at least 7 days of data (i.e., at least 7 kill counts)
        // Alternatively, you might want to check for the number of days with any kills
        let daysWithData = 0;
        for (let day = 0; day < 7; day++) {
            const totalKillsPerDay = data[day].reduce((a, b) => a + b, 0);
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
            xLabels: xLabels,
            yLabels: yLabels,
            backgroundColor: function (context) {
                const value = context.dataset.data[context.dataIndex].v;
                // Define a color scale based on the maximum kill count
                const alpha = maxKills > 0 ? Math.min(value / maxKills, 1) : 0;
                return `rgba(255, 99, 132, ${alpha})`;
            },
        };

        return { labels: [], datasets: [dataset] };
    },
};

export default killsHeatmapChartConfig;
