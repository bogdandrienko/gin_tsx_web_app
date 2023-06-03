// TODO download modules ///////////////////////////////////////////////////////////////////////////////////////////////

import React from "react";

// TODO custom modules /////////////////////////////////////////////////////////////////////////////////////////////////

import axios from "axios";

// TODO export /////////////////////////////////////////////////////////////////////////////////////////////////////////

export default function Page(): JSX.Element {

  async function getData(){
    const response = await axios.get("tasks")
    console.log(response)
  }

    async function sendData(){
      const formData = new FormData()
        formData.append("title", "Amon Ra " + Date.now().toString())
        formData.append("author", "V.Pelevin")
        const response = await axios.post("tasks", formData)
        console.log(response)
    }

  // TODO return ///////////////////////////////////////////////////////////////////////////////////////////////////////
  return (
      <div className="container container-fluid m-5 p-5">
        <button className="btn btn-lg btn-outline-primary m-5" onClick={getData}>getData to console</button>
        <button className="btn btn-lg btn-outline-danger m-5" onClick={sendData}>sendData to backend</button>
      </div>
  );
}
