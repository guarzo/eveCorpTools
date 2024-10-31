// static/js/chartConfigs/valueOverTimeChartConfig.js

import {getCommonOptions, validateChartData} from '../utils.js';

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
    options: getCommonOptions('Isk Destoryed', {
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
        const chartName = 'Isk Destroyed';
        if (!validateChartData(data, chartName)) {
            // Return empty labels and datasets to trigger the noDataPlugin
            return { labels: [], datasets: [] };
        }
        const labels = data.map(item => new Date(item.Time));
        const values = data.map(item => item.Value || 0);

        const datasets = [{
            label: 'Ship Isk Value',
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
