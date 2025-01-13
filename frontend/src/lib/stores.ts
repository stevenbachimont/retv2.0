import { writable } from 'svelte/store';
import type { UserActivity } from './types';

export const userActivities = writable<UserActivity[]>([]);
export const totalCarbon = writable<number>(0);

interface User {
    id: string;
    email: string;
}

// Initialiser explicitement à null et exporter le type
export const user = writable<User | null>(null);
console.log("Store initialized"); 