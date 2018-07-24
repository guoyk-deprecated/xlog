<template>
<div id="app">
<el-container>
  <el-header>
    <span id="brand"><b>XLOG</b></span>
    <el-select v-model="quickPeriod" :disabled="isQuickPeriodPickerDisabled" v-on:change="onQuickPeriodChanged" placeholder="快速选择时间" size="small" style="width:10rem;">
      <el-option v-for="item in quickPeriodOptions" :key="item.value" :label="item.label" :value="item.value"></el-option>
    </el-select>
    <el-date-picker v-model="date" v-if="isDatePickerVisible" :disabled="isDatePickerDisabled" v-on:change="onDateChanged" :clearable="false" type="date" placeholder="选择日期" size="small" style="width: 10rem;"></el-date-picker>
    <el-time-picker v-model="period" v-if="isPeriodPickerVisible" :disabled="isPeriodPickerDisabled" :clearable="false" is-range range-separator="-" start-placeholder="起始时间" end-placeholder="结束时间" size="small" style="width:16rem;"></el-time-picker>
    <el-form v-model="searchForm" id="form-search" :inline="true" size="small" @submit.native.prevent>
     <el-form-item>
        <el-tooltip class="item" effect="dark" content="冒号 : 分割键值对，逗号 , 分割多个值，分号 ; 分割多个键值对" placement="bottom-start">
          <el-input id="input-query" :disabled="isQueryInputDisabled" v-model="searchForm.query" placeholder="查询条件，如 'topic : access, err ; project : sma'"></el-input>
        </el-tooltip>
      </el-form-item>
      <el-form-item>
        <el-button plain type="success" :disabled="isQueryButtonDisabled" @click="onSearchClicked" icon="el-icon-search">查询</el-button>
        <el-button id="btn-auto-refresh" :type="autoRefreshBtnType" :disabled="isAutoRefreshButtonDisabled" @click="onAutoRefreshClicked" size="mini" icon="el-icon-refresh" circle></el-button>
      </el-form-item>
    </el-form>
  </el-header>
  <el-container>
    <el-main>
      <el-row style="margin-bottom: 16px;">
        <el-col>
          <trends-chart :height="50" :chart-data="trendsChartData"></trends-chart>
        </el-col>
      </el-row>
      <el-row v-for="record in records" :key="record._id">
       <el-col>
          <div class="record-meta">
            <small><i class="el-icon-info"></i></small>
            <span class="meta-field">timestamp:</span><span>{{record.timestamp}}</span><span class="meta-semicolon">;</span>
            <span class="meta-field">crid:</span><span>{{record.crid}}</span><span class="meta-semicolon">;</span>
            <span class="meta-field">hostname:</span><span>{{record.hostname}}</span><span class="meta-semicolon">;</span>
            <span class="meta-field">env:</span><span>{{record.env}}</span><span class="meta-semicolon">;</span>
            <span class="meta-field">project:</span><span>{{record.project}}</span><span class="meta-semicolon">;</span>
            <span class="meta-field">topic:</span><span>{{record.topic}}</span>
          </div>
          <pre class="record-message">{{record.message}}</pre>
        </el-col>
      </el-row>
      </el-main>
  </el-container>
</el-container>
</div>
</template>

<script>
import TrendsChart from "./components/TrendsChart.js";
import moment from "moment";
import VueTimers from "vue-timers/mixin";
import { decodeXQL, convertTrendsChartData, formatMoment } from "./utils";

// constants

const quickOptions = [
  {
    label: "手动选择时间",
    value: "manual"
  },
  {
    label: "全天",
    value: "allday"
  },
  {
    label: "最近 15 分钟",
    value: "15m"
  },
  {
    label: "最近 30 分钟",
    value: "30m"
  },
  {
    label: "最近  1 小时",
    value: "1h"
  },
  {
    label: "最近  6 小时",
    value: "6h"
  }
];

// helper functions
function initialPeriod() {
  let e = moment();
  let b;
  if (e.hour() > 1) {
    b = e.clone().subtract(1, "hour");
  } else {
    b = e
      .clone()
      .hour(0)
      .minute(0)
      .second(0);
  }
  return [b, e];
}

export default {
  components: { TrendsChart },
  mixins: [VueTimers],
  timers: {
    autoRefresh: { time: 10000, repeat: true }
  },
  data() {
    return {
      loadingCounter: 0,
      date: new Date(),
      period: initialPeriod(),
      searchForm: { query: "" },
      quickPeriod: "15m",
      quickPeriodOptions: quickOptions,
      records: [],
      trends: []
    };
  },
  computed: {
    autoRefreshBtnType() {
      return this.timers.autoRefresh.isRunning ? "success" : "";
    },
    trendsChartData() {
      return convertTrendsChartData(this.trends);
    },
    isQuickPeriodPickerDisabled() {
      return this.isLoading() || this.timers.autoRefresh.isRunning;
    },
    isDatePickerVisible() {
      return this.quickPeriod == "manual" || this.quickPeriod == "allday";
    },
    isDatePickerDisabled() {
      return this.isLoading() || this.timers.autoRefresh.isRunning;
    },
    isPeriodPickerVisible() {
      return this.quickPeriod == "manual";
    },
    isPeriodPickerDisabled() {
      return this.isLoading() || this.timers.autoRefresh.isRunning;
    },
    isQueryInputDisabled() {
      return this.isLoading() || this.timers.autoRefresh.isRunning;
    },
    isQueryButtonDisabled() {
      return this.isLoading() || this.timers.autoRefresh.isRunning;
    },
    isAutoRefreshButtonDisabled() {
      return this.isLoading() && !this.timers.autoRefresh.isRunning;
    }
  },
  mounted() {
    this.executeSearch();
    this.executeTrends();
  },
  methods: {
    startLoading() {
      this.loadingCounter++;
    },
    endLoading() {
      this.loadingCounter--;
    },
    isLoading() {
      return this.loadingCounter > 0;
    },
    executeSearch() {
      let q = this.buildQuery();
      if (!q) return;
      this.startLoading();
      this.$http
        .post("/api/search", q)
        .then(res => res.json())
        .then(data => {
          this.endLoading();
          this.records = data.result.records;
        })
        .catch(res => {
          this.endLoading();
          this.$notify({
            title: "API 错误",
            message: res.bodyText.trim(),
            type: "error"
          });
        });
    },
    executeTrends() {
      let q = this.buildQuery();
      if (!q) return;
      this.startLoading();
      this.$http
        .post("/api/trends", q)
        .then(res => res.json())
        .then(data => {
          this.endLoading();
          this.trends = data.trends;
        })
        .catch(res => {
          this.endLoading();
          this.$notify({
            title: "API 错误",
            message: res.bodyText.trim(),
            type: "error"
          });
        });
    },
    onSearchClicked() {
      this.executeSearch();
      this.executeTrends();
    },
    onDateChanged() {
      if (this.quickPeriod == "allday") {
        this.executeSearch();
        this.executeTrends();
      }
    },
    onQuickPeriodChanged(v) {
      switch (v) {
        case "manual": {
          break;
        }
        case "allday": {
          this.executeSearch();
          this.executeTrends();
          break;
        }
        default: {
          this.date = new Date();
          this.executeSearch();
          this.executeTrends();
        }
      }
    },
    onAutoRefreshClicked() {
      if (this.timers.autoRefresh.isRunning) {
        this.$timer.stop("autoRefresh");
      } else {
        this.$timer.start("autoRefresh");
      }
    },
    autoRefresh() {
      this.executeSearch();
      this.executeTrends();
    },
    buildQuery() {
      let q = {};
      // check quick period
      switch (this.quickPeriod) {
        case "manual": {
          let b1 = moment(this.period[0]);
          let e1 = moment(this.period[1]);
          let b = moment(this.date)
            .hour(b1.hour())
            .minute(b1.minute())
            .second(b1.second());
          let e = moment(this.date)
            .hour(e1.hour())
            .minute(e1.minute())
            .second(e1.second());
          q.timestamp = {
            beginning: formatMoment(b),
            end: formatMoment(e),
            ascendant: true
          };
          break;
        }
        case "allday": {
          // all day means beginning of day and end of day
          q.timestamp = {
            beginning: formatMoment(moment(this.date).startOf("day")),
            end: formatMoment(moment(this.date).endOf("day"))
          };
          break;
        }
        default: {
          // end is current time
          let end = moment();
          // beginning is a clone of end
          let beginning = end.clone();
          // modify beginning time
          if (this.quickPeriod.endsWith("m")) {
            let num = Number(
              this.quickPeriod.substring(0, this.quickPeriod.length - 1)
            );
            beginning.subtract(num, "minutes");
          } else if (this.quickPeriod.endsWith("h")) {
            let num = Number(
              this.quickPeriod.substring(0, this.quickPeriod.length - 1)
            );
            beginning.subtract(num, "hours");
          } else {
            throw new Error("impossible value");
          }
          // adjust beginning time if exceeded
          if (beginning.day() != end.day()) {
            beginning = end
              .clone()
              .hour(0)
              .minute(0)
              .second(0);
          }
          // assign q.timestamp
          q.timestamp = {
            beginning: formatMoment(beginning),
            end: formatMoment(end)
          };
        }
      }
      decodeXQL(this.searchForm.query, q);
      return q;
    }
  }
};
</script>

<style>
body {
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
    "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
  padding: 0;
  margin: 0;
}

.el-header {
  background-color: #e9eef3;
  padding-top: 16px;
}

.el-main {
  color: #333;
}

#brand {
  margin-right: 0.8rem;
}

#form-search {
  float: right;
}

#input-query {
  width: 24rem;
}

@media screen and (min-width: 1600px) {
  #input-query {
    width: 40rem;
  }
}

#btn-auto-refresh {
  width: 32px;
  height: 32px;
  padding: 8px;
}

.record-meta {
  font-family: monospace;
  padding-left: 8px;
  padding-right: 8px;
}

.meta-field {
  color: #d35400;
}

.meta-semicolon {
  color: #d35400;
}

pre.record-message {
  border-radius: 2px;
  padding: 8px;
  background-color: #f8f8f8;
  word-wrap: break-word;
  white-space: pre-wrap;
}
</style>
