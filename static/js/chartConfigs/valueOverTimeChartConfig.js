// chartConfigs/valueOverTimeChartConfig.js
import { commonOptions } from '../utils.js';

const valueOverTimeChartConfig = {
    id: 'valueOverTimeChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdValueOverTimeData',
        ytd: 'ytdValueOverTimeData',
        lastMonth: 'lastMValueOverTimeData',
    },
    type: 'line',
    options: {
        ...commonOptions,
        plugins: {
            ...commonOptions.plugins,
            legend: { display: false },
        },
        scales: {
            x: {
                ...commonOptions.scales.x,
                type: 'time',
                time: {
                    unit: 'day',
                },
            },
            y: {
                ...commonOptions.scales.y,
            },
        },
    },
    processData: function (data) {
        // ... processing logic remains the same ...
        const labels = data.map(item => new Date(item.Time));
        const values = data.map(item => item.Value || 0);

        const datasets = [{
            label: 'ISK Value Destroyed Over Time',
            data: values,
            borderColor: 'rgba(54, 162, 235, 1)',
            backgroundColor: 'rgba(54, 162, 235, 0.5)',
            fill: true,
        }];

        return { labels, datasets };
    },
};

export default valueOverTimeChartConfig;
