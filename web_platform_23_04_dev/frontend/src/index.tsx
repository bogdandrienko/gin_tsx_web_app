// TODO download modules ///////////////////////////////////////////////////////////////////////////////////////////////

import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import axios from "axios";

// TODO custom modules /////////////////////////////////////////////////////////////////////////////////////////////////

import HomePage from "./data/pages/HomePage";
import './data/css/bootstrap/bootstrap.min.css';
import './data/css/font_awesome/css/all.css';
import './data/css/font_awesome/css/all.css';

// TODO settings ///////////////////////////////////////////////////////////////////////////////////////////////////////

axios.defaults.baseURL = "http://127.0.0.1:8080/api/";
axios.defaults.headers.common = {
    ...axios.defaults.headers.common,
    "Access-Control-Allow-Origin": "*",
};

createRoot(document.getElementById("root")!).render(
  // <React.StrictMode>
      <BrowserRouter>
          <Routes>
              <Route path="/" element={<HomePage />}></Route>
          </Routes>
      </BrowserRouter>
  // </React.StrictMode>
);
