// static/js/chartConfigs/killActivityChartConfig.js

import { getCommonOptions, validateChartData } from '../utils.js';

/**
 * Configuration for the Kill Activity Over Time Chart
 */
const killActivityChartConfig = {
    id: 'killActivityChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdKillActivityData',
        ytd: 'ytdKillActivityData',
        lastMonth: 'lastMKillActivityData',
    },
    type: 'line',
    options: getCommonOptions('Kill Activity Over Time', {
        plugins: {
            legend: { display: true, position: 'top', labels: { color: '#ffffff' } },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const label = context.dataset.label || '';
                        const value = context.parsed.y !== undefined ? context.parsed.y : context.parsed.x;
                        return `${label}: ${value}`;
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day',
                    tooltipFormat: 'MMM D, YYYY', // Customize tooltip date format
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
    }),
    processData: function (data) {
        const chartName = 'Kill Activity Over Time Chart';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }

        // Assuming 'Time' is a date string and 'Kills' is the number of kills
        const labels = data.map(item => new Date(item.Time));
        const kills = data.map(item => item.Kills || 0);

        const datasets = [{
            label: 'Kills Over Time',
            data: kills,
            borderColor: 'rgba(255, 77, 77, 1)',
            backgroundColor: 'rgba(255, 77, 77, 0.5)',
            fill: true,
            tension: 0.4, // Smooth the line
            pointBackgroundColor: 'rgba(255, 77, 77, 1)',
            pointBorderColor: '#fff',
            pointHoverBackgroundColor: '#fff',
            pointHoverBorderColor: 'rgba(255, 77, 77, 1)',
        }];

        return { labels, datasets };
    },
};

export default killActivityChartConfig;
