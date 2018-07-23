import moment from "moment";

const XQLFields = ["crid", "env", "project", "topic", "hostname"];

export function decodeXQL(s, out) {
  let kvs = s.split(";");
  for (let kv of kvs) {
    kv = kv.trim();
    let knv = kv.split(":");
    if (knv.length == 2) {
      let k = knv[0];
      let v = knv[1];
      k = k.trim();
      v = v.trim();
      if (v.length == 0) continue;
      if (XQLFields.includes(k)) {
        let vs = v.split(",");
        if (vs.length > 1) {
          let nvs = [];
          for (let sv of vs) {
            nvs.push(sv.trim());
          }
          out[k] = nvs.join(",");
        } else {
          out[k] = v;
        }
      }
    }
  }
}

export function convertTrendsChartData(trends) {
  let data = {
    labels: [],
    datasets: [
      {
        data: [],
        borderWidth: 1,
        backgroundColor: "#7f8c8d"
      }
    ]
  };
  for (let trend of trends) {
    let date = new Date();
    data.labels.push(
      moment(trend.beginning)
        .utc()
        .format("HH:mm")
    );
    data.datasets[0].data.push(trend.count);
  }
  return data;
}

export function formatMoment(mm) {
  return mm.format("YYYY-MM-DDTHH:mm:ss[Z]");
}
