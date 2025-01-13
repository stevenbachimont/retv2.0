export interface CarbonActivity {
    id: string;
    category: 'transport' | 'alimentation' | 'logement' | 'consommation';
    name: string;
    carbonFactor: number; // en kg CO2e
    unit: string;
}

export interface UserActivity {
    activityId: string;
    quantity: number;
    date: Date;
}

export type CarbonData = {
    Transports: any;
    Logement_electromenagers: any;
    Alimentation: any;
    Vetements: any;
} 