// static/js/chartConfigs/damageFinalBlowsChartConfig.js

import { truncateLabel, getColor, getCommonOptions } from '../utils.js';

/**
 * Configuration for the Damage Done and Final Blows Chart
 */
const damageFinalBlowsChartConfig = {
    id: 'damageFinalBlowsChart',
    instance: null,
    dataKeys: {
        mtd: 'mtdCharacterDamageData',
        ytd: 'ytdCharacterDamageData',
        lastMonth: 'lastMCharacterDamageData',
    },
    type: 'bar',
    options: getCommonOptions('Damage Done and Final Blows'),
    processData: function (data) {
        const labels = data.map(item => item.Name || 'Unknown');
        const damageData = data.map(item => item.DamageDone || 0);
        const finalBlowsData = data.map(item => item.FinalBlows || 0);

        const fullLabels = [...labels];
        const truncatedLabels = labels.map(label => truncateLabel(label, 10));

        const datasets = [
            {
                label: 'Damage Done',
                data: damageData,
                backgroundColor: 'rgba(255, 77, 77, 0.7)',
            },
            {
                label: 'Final Blows',
                data: finalBlowsData,
                backgroundColor: 'rgba(54, 162, 235, 0.7)',
            },
        ];

        return { labels: truncatedLabels, datasets, fullLabels };
    },
};

export default damageFinalBlowsChartConfig;
