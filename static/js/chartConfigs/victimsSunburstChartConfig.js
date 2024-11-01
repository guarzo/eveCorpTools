// static/js/chartConfigs/victimsSunburstChartConfig.js

import { getColor, getCommonOptions, validateChartDataArray } from '../utils.js';

/**
 * Configuration for the Victims Sunburst Chart
 */
const victimsSunburstChartConfig = {
    id: 'victimsSunburstChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdVictimsSunburstData',
        ytd: 'ytdVictimsSunburstData',
        lastMonth: 'lastMVictimsSunburstData',
    },
    type: 'sunburst',
    options: getCommonOptions('Victims Sunburst', {
        plugins: {
            legend: { display: false },
            tooltip: {
                callbacks: {
                    label: function (context) {
                        const label = context.raw.name || '';
                        const value = context.raw.value || 0;
                        return `${label}: ${value}`;
                    },
                },
            },
        },
    }),
    processData: function (data) {
        const chartName = 'Victims Sunburst Chart';
        if (!validateChartDataArray(data, chartName)) {
            return { datasets: [] };
        }

        return {
            datasets: [{
                data: data,
                backgroundColor: function (context) {
                    const index = context.dataIndex;
                    return getColor(index);
                },
                borderWidth: 1,
            }],
        };
    },
};

export default victimsSunburstChartConfig;
