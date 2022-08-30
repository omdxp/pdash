import { Component, createSignal } from "solid-js";
import { Link, useNavigate } from "@solidjs/router";

import { signup } from "../../store";

const Signup: Component = () => {
  const navigate = useNavigate();

  const [fullname, setFullname] = createSignal("");
  const [username, setUsername] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [confirmPassword, setConfirmPassword] = createSignal("");
  return (
    <section class="bg-gray-50 dark:bg-gray-900">
      <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
        <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
          <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
            <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
              Create an account
            </h1>
            <form class="space-y-4 md:space-y-6" action="#">
              <div>
                <label
                  for="fullname"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Your fullname
                </label>
                <input
                  type="text"
                  name="fullname"
                  id="fullname"
                  class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="cool-user"
                  value={fullname()}
                  onChange={(e) =>
                    setFullname((e.target as HTMLInputElement).value)
                  }
                  required
                />
              </div>
              <div>
                <label
                  for="username"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Your username
                </label>
                <input
                  type="text"
                  name="username"
                  id="username"
                  class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="cool-user"
                  value={username()}
                  onChange={(e) =>
                    setUsername((e.target as HTMLInputElement).value)
                  }
                  required
                />
              </div>
              <div>
                <label
                  for="email"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Your email
                </label>
                <input
                  type="email"
                  name="email"
                  id="email"
                  class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  placeholder="name@company.com"
                  value={email()}
                  onChange={(e) =>
                    setEmail((e.target as HTMLInputElement).value)
                  }
                  required
                />
              </div>
              <div>
                <label
                  for="password"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Password
                </label>
                <input
                  type="password"
                  name="password"
                  id="password"
                  placeholder="••••••••"
                  class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  value={password()}
                  onChange={(e) =>
                    setPassword((e.target as HTMLInputElement).value)
                  }
                  required
                />
              </div>
              <div>
                <label
                  for="confirm-password"
                  class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                >
                  Confirm password
                </label>
                <input
                  type="password"
                  name="confirm-password"
                  id="confirm-password"
                  placeholder="••••••••"
                  class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-primary-600 focus:border-primary-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                  value={confirmPassword()}
                  onChange={(e) =>
                    setConfirmPassword((e.target as HTMLInputElement).value)
                  }
                  required
                />
              </div>
              <button
                type="submit"
                onClick={async (e) => {
                  e.preventDefault();
                  if (
                    fullname() === "" ||
                    username() === "" ||
                    email() === "" ||
                    password() === ""
                  ) {
                    alert("Please fill out all fields");
                    return;
                  }
                  if (password() !== confirmPassword()) {
                    alert("Passwords do not match");
                    return;
                  }
                  const res = await signup({
                    fullname: fullname(),
                    username: username(),
                    email: email(),
                    password: password(),
                  });
                  if (res) {
                    navigate("/login", { replace: true });
                  }
                }}
                class="w-full text-white bg-slate-600 hover:bg-slate-500 focus:ring-4 focus:outline-none focus:ring-primary-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-primary-600 dark:hover:bg-primary-700 dark:focus:ring-primary-800"
              >
                Create an account
              </button>
            </form>
            <p class="text-sm font-light text-gray-500 dark:text-gray-400">
              Already have an account?{" "}
              <Link href="/login">
                <a class="font-medium text-primary-600 hover:underline dark:text-primary-500">
                  Login here
                </a>
              </Link>
            </p>
          </div>
        </div>
      </div>
    </section>
  );
};
export default Signup;
