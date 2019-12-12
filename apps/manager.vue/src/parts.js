import home      from './home/routes';
import workspace from '../../workspace/data_vue/routes';
import flow      from '../../workspace/flow_vue/routes';

let routes = [
  ...home,
  ...workspace,
  ...flow,
];

export default routes;
