// static/js/chartConfigs/ourLossesCombinedChartConfig.js

import { truncateLabel, getCommonOptions } from '../utils.js';

/**
 * Configuration for the Our Losses Combined Chart
 */
const ourLossesCombinedChartConfig = {
    id: 'ourLossesCombinedChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdOurLossesValueData',
        ytd: 'ytdOurLossesValueData',
        lastMonth: 'lastMOurLossesValueData',
    },
    type: 'bar',
    options: getCommonOptions('Our Losses Combined', {
        plugins: {
            legend: { display: true },
            tooltip: {
                mode: 'index',
                intersect: false,
                callbacks: {
                    title: function (context) {
                        const index = context[0].dataIndex;
                        return context[0].chart.data.fullLabels[index];
                    },
                },
            },
        },
        scales: {
            x: {
                type: 'category',
                ticks: {
                    color: '#ffffff',
                    maxRotation: 45,
                    minRotation: 45,
                    autoSkip: false,
                },
                grid: { display: false },
            },
            y: {
                type: 'linear',
                position: 'left',
                ticks: { color: '#ffffff' },
                grid: { color: '#444' },
                beginAtZero: true,
            },
            y1: {
                type: 'linear',
                position: 'right',
                ticks: { color: '#ffffff' },
                grid: { drawOnChartArea: false },
                beginAtZero: true,
            },
        },
    }),
    processData: function (data) {
        // Filter out any invalid data entries
        data = data.filter(item => item && item.CharacterName);

        const labels = data.map(item => item.CharacterName);
        const lossesValueData = data.map(item => item.TotalValue || 0);
        const lossesCountData = data.map(item => item.LossesCount || 0);
        const shipCountData = data.map(item => item.ShipCount || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const datasets = [
            {
                label: 'Losses Value (ISK)',
                data: lossesValueData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
                yAxisID: 'y1',
            },
            {
                label: 'Losses Count',
                data: lossesCountData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
                yAxisID: 'y',
            },
            {
                label: 'Ship Types Lost',
                data: shipCountData,
                backgroundColor: 'rgba(75, 192, 192, 0.7)',
                yAxisID: 'y',
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels };
    },
};

export default ourLossesCombinedChartConfig;
