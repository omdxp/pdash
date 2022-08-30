import { Link, useNavigate } from "@solidjs/router";
import { setToken, token } from "../../store";

import { Component } from "solid-js";

const Header: Component = () => {
  const navigate = useNavigate();
  return (
    <header>
      <nav class="bg-slate-800 border-gray-200 px-2 sm:px-4 py-2.5  dark:bg-gray-900">
        <div class="container flex flex-wrap justify-between items-center mx-auto">
          {token() ? (
            <>
              <Link href="/">
                <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
                  Home
                </a>
              </Link>
              <Link href="/suppliers">
                <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
                  Suppliers
                </a>
              </Link>
              <Link href="/customers">
                <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
                  Customers
                </a>
              </Link>
              <button
                onClick={(e) => {
                  e.preventDefault();
                  setToken(null);
                  navigate("/login", { replace: true });
                }}
                class="text-gray-300 font-semibold text-lg hover:text-gray-600"
              >
                Logout
              </button>
            </>
          ) : (
            <Link href="/login">
              <a class="text-gray-300 font-semibold text-lg hover:text-gray-600">
                Login
              </a>
            </Link>
          )}
        </div>
      </nav>
    </header>
  );
};
export default Header;
