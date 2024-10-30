// static/js/chartConfigs/valueOverTimeChartConfig.js

import { getCommonOptions } from '../utils.js';

/**
 * Configuration for the Value Over Time Chart
 */
const valueOverTimeChartConfig = {
    id: 'valueOverTimeChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdValueOverTimeData',
        ytd: 'ytdValueOverTimeData',
        lastMonth: 'lastMValueOverTimeData',
    },
    type: 'line',
    options: getCommonOptions('Value Over Time', {
        plugins: {
            legend: { display: false },
        },
        scales: {
            x: {
                type: 'time',
                time: {
                    unit: 'day',
                },
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
            },
            y: {
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
            },
        },
    }),
    processData: function (data) {
        const labels = data.map(item => new Date(item.Time));
        const values = data.map(item => item.Value || 0);

        const datasets = [{
            label: 'ISK Value Destroyed Over Time',
            data: values,
            borderColor: 'rgba(54, 162, 235, 1)',
            backgroundColor: 'rgba(54, 162, 235, 0.5)',
            fill: true,
            tension: 0.4, // Smooth the line
        }];

        return { labels, datasets };
    },
};

export default valueOverTimeChartConfig;
