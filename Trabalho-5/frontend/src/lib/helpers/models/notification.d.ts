import type { Lance } from './lance';
import type { Status } from './status';
import type { Link } from './link';

export interface Notification {
	type: string;
	lance?: Lance;
	statusData?: Status;
	linkData?: Link;
}
