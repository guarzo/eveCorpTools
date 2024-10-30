// chartConfigs/victimsSunburstChartConfig.js
import { getColor, commonOptions } from '../utils.js';

const victimsSunburstChartConfig = {
    id: 'victimsSunburstChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdVictimsSunburstData',
        ytd: 'ytdVictimsSunburstData',
        lastMonth: 'lastMVictimsSunburstData',
    },
    type: 'sunburst',
    options: {
        ...commonOptions,
        plugins: {
            ...commonOptions.plugins,
            legend: { display: false },
            tooltip: {
                ...commonOptions.plugins.tooltip,
                callbacks: {
                    label: function (context) {
                        const label = context.raw.name || '';
                        const value = context.raw.value || 0;
                        return `${label}: ${value}`;
                    },
                },
            },
        },
    },
    processData: function (data) {
        const datasets = [{
            data: data,
            backgroundColor: function (context) {
                const index = context.dataIndex;
                return getColor(index);
            },
        }];

        return { datasets };
    },
};

export default victimsSunburstChartConfig;
