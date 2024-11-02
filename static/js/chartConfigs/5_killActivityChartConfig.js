// static/js/chartConfigs/killActivityChartConfig.js

import { getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Kill Activity Over Time Chart
 */
const killActivityChartConfig = {
    id: 'killActivityChart',
    instance: {}, // Initialize as an object to store chart instances per timeframe
    dataKeys: {
        mtd: { dataVar: 'mtdKillActivityData', canvasId: 'killActivityChart_mtd' },
        ytd: { dataVar: 'ytdKillActivityData', canvasId: 'killActivityChart_ytd' },
        lastMonth: { dataVar: 'lastMKillActivityData', canvasId: 'killActivityChart_lastM' },
    },
    type: 'line',
    dataType: 'array', // Specify that this chart expects array data
    options: getCommonOptions('Kill Activity Over Time', {
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                mode: 'nearest', // Changed for better tooltip behavior
                intersect: true,
                callbacks: {
                    label: function (context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed.x;
                        return `${label}: ${value}`;
                    },
                },
            },
            datalabels: {
                color: '#ffffff',
                align: 'top',
                formatter: (value) => `${value.y}`, // Access y value from data point
                font: {
                    size: 10,
                    weight: 'bold',
                },
            },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day', // Adjust based on your data granularity
                    tooltipFormat: 'MMM d, yyyy', // Customize tooltip date format
                    displayFormats: {
                        day: 'MMM d',
                        hour: 'MMM d, hA',
                    },
                },
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Time',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
            y: {
                beginAtZero: true,
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                title: {
                    display: true,
                    text: 'Kills',
                    color: '#ffffff',
                    font: {
                        size: 14,
                        family: 'Montserrat, sans-serif',
                        weight: 'bold',
                    },
                },
            },
        },
        responsive: true,
        maintainAspectRatio: false, // Allow the chart to adjust its height
    }),
    processData: function (data) {
        const chartName = 'Kill Activity Over Time Chart';
        if (!validateChartDataArray(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [], noDataMessage: 'No data available for this chart.' };
        }

        // Map data to {x: Time, y: Kills} format
        const dataPoints = data.map(item => {
            const date = new Date(item.Time);
            const kills = item.Kills || 0;

            if (isNaN(date)) {
                console.warn(`Invalid date format in data for ${chartName}:`, item.Time);
                return null; // Exclude invalid data points
            }

            return { x: date, y: kills };
        }).filter(point => point !== null); // Remove nulls

        // Check if there are at least 7 days of data
        if (dataPoints.length < 3) {
            console.warn(`Not enough data points (${dataPoints.length}) for ${chartName}.`);
            return { labels: [], datasets: [], noDataMessage: 'Not enough data to display the chart.' };
        }

        const datasets = [{
            label: 'Kills Over Time',
            data: dataPoints,
            borderColor: 'rgba(255, 77, 77, 1)',
            backgroundColor: 'rgba(255, 77, 77, 0.5)',
            fill: true,
            tension: 0.4, // Smooth the line
            pointBackgroundColor: 'rgba(255, 77, 77, 1)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(255, 77, 77, 1)',
        }];

        return { labels: [], datasets };
    },
};

export default killActivityChartConfig;
