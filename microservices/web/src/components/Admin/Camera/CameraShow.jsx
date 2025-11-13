import { Show, SimpleShowLayout } from "react-admin";
import CameraStream from "./Camera";

const CameraShow = () => (
  <Show>
    <SimpleShowLayout>
      <CameraStream />
    </SimpleShowLayout>
  </Show>
);

export default CameraShow;
