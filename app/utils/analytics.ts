import mixpanel, { Dict } from 'mixpanel-browser';

const mixPanelId = process.env.NEXT_PUBLIC_MIXPANEL_ID;

export const init = () => {
  if (!mixPanelId?.length) {
    return;
  }

  mixpanel.init(mixPanelId);
};

export const track = (eventName: string, props?: Dict) => {
  if (!mixPanelId?.length) {
    return;
  }
  mixpanel.track(eventName, props);
};
