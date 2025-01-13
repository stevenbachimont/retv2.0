import { writable } from 'svelte/store';
import type { UserActivity } from './types';

export const userActivities = writable<UserActivity[]>([]);
export const totalCarbon = writable<number>(0); 