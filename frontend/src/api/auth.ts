type Register = {
  email: string;
  name: string;
  phone: string;
  password: string;
};

export async function register(registerData: Register) {
  const res = await fetch('/api/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(registerData),
    // VERY IMPORTANT: sends cookies cross-origin
    credentials: 'include',
  });

  if (!res.ok) throw new Error(await res.text());
  return await res.json();
}

export async function login(data: { email: string; password: string }) {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
    credentials: 'include',
  });

  if (!res.ok) throw new Error(await res.text());
  return await res.json();
}

export async function fetchCurrentUser(): Promise<User> {
  const res = await fetch('/api/me', { credentials: 'include' });
  if (!res.ok) throw new Error('Not authenticated');
  return res.json();
}

export async function logout() {
  const res = await fetch('/api/logout', {
    method: 'POST',
    credentials: 'include',
  });
  if (!res.ok) throw new Error(await res.text());
  return;
}
