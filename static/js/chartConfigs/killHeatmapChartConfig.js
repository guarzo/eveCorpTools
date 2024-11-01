// static/js/chartConfigs/killsHeatmapChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kills Heatmap Chart
 */
const killsHeatmapChartConfig = {
    id: 'killsHeatmapChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillHeatmapData',
        ytd: 'ytdKillHeatmapData',
        lastMonth: 'lastMKillHeatmapData',
    },
    type: 'matrix', // Using matrix chart type from chartjs-chart-matrix
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
            return { labels: [], datasets: [] };
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

        for (let day = 0; day < 7; day++) {
            for (let hour = 0; hour < 24; hour++) {
                const kills = data[day][hour] || 0;
                matrixData.push({
                    x: hour.toString(),
                    y: yLabels[day],
                    v: kills,
                });
            }
        }

        const dataset = {
            label: 'Kills Heatmap',
            data: matrixData,
            xLabels: xLabels,
            yLabels: yLabels,
            backgroundColor: function (context) {
                const value = context.dataset.data[context.dataIndex].v;
                // Define a color scale (e.g., higher kills = darker color)
                const alpha = value > 0 ? Math.min(value / 100, 1) : 0; // Adjust scaling based on data
                return `rgba(255, 99, 132, ${alpha})`;
            },
        };

        return { labels: [], datasets: [dataset] };
    },
};

export default killsHeatmapChartConfig;
