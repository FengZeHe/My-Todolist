<template>
  <div class="main">
    <el-container>
      <el-header class="todo-header">
        <p> To do list</p>
        <el-input
          v-model="input"
          placeholder="请输入内容"
          @keyup.enter.native="AddTask"
        >
        </el-input>
      </el-header>
      
      <el-main>
        <transition name="el-fade-in-linear">
          <div v-if="todotask.length">
            <ol>
              <li v-for="event in todotask" :key="event.TaskID">
                <div class="container">
                  {{ event.task_name }}
                  <div class="content">
                    <el-button
                      type="success"
                      icon="el-icon-check"
                      circle
                      @click="HandleTask(event.TaskID, !event.task_state)"
                    ></el-button>
                    <el-button
                      type="danger"
                      icon="el-icon-delete"
                      circle
                      @click="DropTask(event.TaskID)"
                    ></el-button>
                  </div>
                </div>
                <el-divider></el-divider>
              </li>
            </ol>
          </div>
        </transition>
        <!--  已完成区域 -->
        <transition name="el-fade-in-linear">
          <div v-if="donetask.length">
            <p>已完成</p>
            <ol>
              <li v-for="event in donetask" :key="event.TaskID">
                <div class="container">
                  <del> {{ event.task_name }}</del>
                  <div class="content">
                    <el-button
                      type="primary"
                      icon="el-icon-refresh-left"
                      circle
                      @click="HandleTask(event.TaskID, !event.task_state)"
                    ></el-button>
                    <el-button
                      type="danger"
                      icon="el-icon-delete"
                      circle
                      @click="DropTask(event.TaskID)"
                    ></el-button>
                  </div>
                </div>
                <el-divider></el-divider>
              </li>
            </ol>
          </div>
        </transition>
      </el-main>
    </el-container>
  </div>
</template>

<script>
export default {
  name: "MainPage",
  data() {
    return {
      todotask: [],
      donetask: [],
      hasbook: true,
      info: null,
      input: "",
    };
  },
  methods: {
    incress() {
      this.count++;
    },
    getTasklist() {
      this.$axios({
        methods: "get",
        url: "http://192.168.2.36:9002/tasklist",
      })
        .then((response) => {
          this.todotask = [];
          this.donetask = [];
          for (var i = 0; i < response.data.length; i++) {
            if (response.data[i].task_state == false) {
              this.todotask.push(response.data[i]);
            } else {
              this.donetask.push(response.data[i]);
            }
          }
        })
        .catch((error) => console.log(error, "error"));
    },
    AddTask() {
      this.$axios
        .post(
          "http://192.168.2.36:9002/tasklist",
          JSON.stringify({ task_name: this.input, task_detail: " " })
        )
        .then((response) => {
          console.log(response);
          this.getTasklist();
          this.input = " ";
        })
        .catch((error) => {
          console.log(error);
        });
    },
    DropTask(TaskID) {
      this.$axios
        .post(
          "http://192.168.2.36:9002/tasklist/delete",
          JSON.stringify({ TaskID: TaskID })
        )
        .then((response) => {
          console.log(response);
          this.getTasklist();
        })
        .catch((error) => {
          console.log(error);
        });
    },
    HandleTask(taskid, StrState) {
      console.log(taskid, StrState);
      var TaskState;
      if (StrState.toString() == "false") {
        TaskState = "0";
      } else {
        TaskState = "1";
      }
      this.$axios
        .post(
          "http://192.168.2.36:9002/tasklist/state",
          JSON.stringify({ TaskID: String(taskid), task_state: TaskState })
        )
        .then((response) => {
          console.log(response);
          this.getTasklist();
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
  mounted() {
    this.getTasklist();
  },
};
</script>

<style scoped>
.main .todo-header{
  height: 100%;
  box-sizing: content-box;
  padding-bottom:50px;
}

li {
  display: flex;
  flex-direction: column;
}
.container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>