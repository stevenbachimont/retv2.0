import { writable } from 'svelte/store';
import type { UserActivity } from './types';

export const userActivities = writable<UserActivity[]>([]);
export const totalCarbon = writable<number>(0);

type User = {
    id: string;
    email: string;
    username: string;
} | null;

// Initialiser explicitement à null et exporter le type
export const user = writable<User | null>(null);
console.log("Store initialized"); 