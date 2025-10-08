import React, { useState } from 'react';
import './App.css';

// API base URL from environment variable or fallback to localhost
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  createdAt: string;
}

interface Task {
  id: number;
  title: string;
  subtasks: Subtask[];
  isExpanded: boolean;
  createdAt: string;
  completedAt?: string;
}

interface Subtask {
  id: number;
  title: string;
  completed: boolean;
}

function App() {
  const [task, setTask] = useState('');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [user, setUser] = useState<User | null>(null);
  const [showProfile, setShowProfile] = useState(false);
  const [showHistory, setShowHistory] = useState(false);
  const [taskHistory, setTaskHistory] = useState<Task[]>([]);
  const [showRegister, setShowRegister] = useState(false);
  const [showLogin, setShowLogin] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [registerForm, setRegisterForm] = useState({
    username: '',
    email: '',
    password: '',
    firstName: '',
    lastName: ''
  });
  const [loginForm, setLoginForm] = useState({
    username: '',
    password: ''
  });

  // Инициализация пользователя и загрузка данных
  React.useEffect(() => {
    const savedUser = localStorage.getItem('taskSplitter_user');
    const savedTasks = localStorage.getItem('taskSplitter_tasks');
    const savedHistory = localStorage.getItem('taskSplitter_history');

    if (savedUser) {
      const userData = JSON.parse(savedUser);
      setUser(userData);
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
    }

    if (savedTasks) {
      setTasks(JSON.parse(savedTasks));
    }

    if (savedHistory) {
      setTaskHistory(JSON.parse(savedHistory));
    }
  }, []);

  // Сохранение задач в localStorage
  React.useEffect(() => {
    if (tasks.length > 0) {
      localStorage.setItem('taskSplitter_tasks', JSON.stringify(tasks));
    }
  }, [tasks]);

  // Сохранение истории в localStorage
  React.useEffect(() => {
    if (taskHistory.length > 0) {
      localStorage.setItem('taskSplitter_history', JSON.stringify(taskHistory));
    }
  }, [taskHistory]);

  // Функция входа
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!loginForm.username || !loginForm.password) {
      setError('Заполните все поля');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(loginForm)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка входа');
      }

      const result = await response.json();
      console.log('✅ Пользователь вошел:', result);

      // Создаем пользователя для фронтенда
      const loggedInUser: User = {
        id: result.user.sub || 'user_' + Date.now(),
        name: `${result.user.given_name || ''} ${result.user.family_name || ''}`.trim() || result.user.preferred_username,
        email: result.user.email,
        avatar: '👤',
        createdAt: new Date().toISOString()
      };

      setUser(loggedInUser);
      localStorage.setItem('taskSplitter_user', JSON.stringify(loggedInUser));
      setIsAuthenticated(true);
      setShowLogin(false);
      setLoginForm({ username: '', password: '' });
      setError('');

    } catch (err) {
      console.error('❌ Ошибка входа:', err);
      setError(`Ошибка входа: ${err instanceof Error ? err.message : 'Неизвестная ошибка'}`);
    } finally {
      setIsLoading(false);
    }
  };

  // Функция выхода
  const handleLogout = () => {
    setUser(null);
    setIsAuthenticated(false);
    localStorage.removeItem('taskSplitter_user');
    setTasks([]);
    setTaskHistory([]);
    setError('');
  };

  // Функция регистрации
  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!registerForm.username || !registerForm.email || !registerForm.password) {
      setError('Заполните все обязательные поля');
      return;
    }

    setIsLoading(true);
    setError('');

    try {
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(registerForm)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка регистрации');
      }

      const result = await response.json();
      console.log('✅ Пользователь зарегистрирован:', result);

      // Создаем пользователя для фронтенда
      const newUser: User = {
        id: result.user.id.toString(),
        name: `${result.user.first_name} ${result.user.last_name}`.trim() || result.user.username,
        email: result.user.email,
        avatar: '👤',
        createdAt: result.user.created_at
      };

      setUser(newUser);
      localStorage.setItem('taskSplitter_user', JSON.stringify(newUser));
      setIsAuthenticated(true);
      setShowRegister(false);
      setRegisterForm({ username: '', email: '', password: '', firstName: '', lastName: '' });
      setError('');

    } catch (err) {
      console.error('❌ Ошибка регистрации:', err);
      setError(`Ошибка регистрации: ${err instanceof Error ? err.message : 'Неизвестная ошибка'}`);
    } finally {
      setIsLoading(false);
    }
  };

  // Polling для получения результата разбивки
  const pollForResult = async (requestId: string, taskId: number) => {
    const maxAttempts = 30; // Максимум 30 попыток (30 секунд)
    let attempts = 0;
    
    const poll = async (): Promise<void> => {
      try {
        const response = await fetch(`${API_BASE_URL}/api/v1/split/${requestId}/status`, {
          headers: {
            'Authorization': 'Bearer mock_token_user123'
          }
        });
        
        if (!response.ok) {
          throw new Error('Ошибка получения статуса');
        }
        
        const statusData = await response.json();
        console.log('📊 Статус разбивки:', statusData);
        
        if (statusData.status === 'completed') {
          // Парсим результат и создаем подзадачи
          const subtasks = parseSubtasks(statusData.result);
          
          // Создаем новую задачу с подзадачами
          const newTask: Task = {
            id: taskId,
            title: task,
            subtasks: subtasks,
            isExpanded: true,
            createdAt: new Date().toISOString()
          };
          
          // Добавляем задачу в список
          setTasks(prevTasks => [newTask, ...prevTasks]);
          
          // Очищаем сообщение об ошибке
          setError('');
          
          console.log('✅ Задача успешно разбита на подзадачи:', newTask);
          return;
        } else if (statusData.status === 'failed') {
          setError(`Ошибка разбивки задачи: ${statusData.error || 'Неизвестная ошибка'}`);
          return;
        }
        
        // Если статус pending, продолжаем polling
        attempts++;
        if (attempts < maxAttempts) {
          setTimeout(poll, 1000); // Ждем 1 секунду перед следующей попыткой
        } else {
          setError('Время ожидания истекло. Попробуйте еще раз.');
        }
      } catch (err) {
        console.error('❌ Ошибка polling:', err);
        setError(`Ошибка получения результата: ${err instanceof Error ? err.message : 'Неизвестная ошибка'}`);
      }
    };
    
    await poll();
  };

  // Парсинг подзадач из JSON результата
  const parseSubtasks = (result: string): Subtask[] => {
    try {
      // Извлекаем JSON из markdown блока
      const jsonMatch = result.match(/```json\n([\s\S]*?)\n```/);
      if (jsonMatch) {
        const jsonStr = jsonMatch[1];
        const parsed = JSON.parse(jsonStr);
        
        return parsed.map((item: any, index: number) => ({
          id: index + 1,
          title: item.title || '',
          completed: false
        }));
      }
      
      // Если не удалось извлечь JSON, возвращаем пустой массив
      return [];
    } catch (err) {
      console.error('❌ Ошибка парсинга подзадач:', err);
      return [];
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!task.trim()) return;

    setIsLoading(true);
    setError('');
    
    try {
      // Сначала создаем задачу в базе данных
      const createTaskResponse = await fetch(`${API_BASE_URL}/api/v1/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer mock_token_user123'
        },
        body: JSON.stringify({
          title: task,
          description: task,
          priority: 'medium'
        })
      });

      if (!createTaskResponse.ok) {
        const errorData = await createTaskResponse.json();
        throw new Error(errorData.error || 'Ошибка создания задачи');
      }

      const createdTask = await createTaskResponse.json();
      console.log('✅ Задача создана:', createdTask);

      // Теперь отправляем задачу на разбивку
      const response = await fetch(`${API_BASE_URL}/api/v1/split`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer mock_token_user123'
        },
        body: JSON.stringify({
          task_id: createdTask.id,
          text: task
        })
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Ошибка сервера');
      }

      const result = await response.json();
      console.log('✅ Запрос отправлен в ChatGPT:', result);
      
      // Показываем сообщение о том, что запрос обрабатывается
      setError('Делаем проще...');
      
      // Polling для получения результата
      await pollForResult(result.request_id, createdTask.id);
      
      // Очищаем поле ввода
      setTask('');
      
    } catch (err) {
      console.error('❌ Ошибка API:', err);
      setError(`Ошибка: ${err instanceof Error ? err.message : 'Неизвестная ошибка'}`);
    } finally {
      setIsLoading(false);
    }
  };

  const toggleSubtask = (taskId: number, subtaskId: number) => {
    setTasks(prev => 
      prev.map(task => 
        task.id === taskId 
          ? {
              ...task,
              subtasks: task.subtasks.map(subtask =>
                subtask.id === subtaskId 
                  ? { ...subtask, completed: !subtask.completed }
                  : subtask
              )
            }
          : task
      )
    );
  };

  const toggleTaskExpansion = (taskId: number) => {
    setTasks(prev => 
      prev.map(task => 
        task.id === taskId 
          ? { ...task, isExpanded: !task.isExpanded }
          : task
      )
    );
  };

  const clearAllTasks = () => {
    // Перемещаем текущие задачи в историю
    const tasksToArchive = tasks.map(task => ({
      ...task,
      completedAt: new Date().toISOString()
    }));
    
    setTaskHistory(prev => [...prev, ...tasksToArchive]);
    setTasks([]);
    setError('');
  };

  const archiveTask = (taskId: number) => {
    const taskToArchive = tasks.find(task => task.id === taskId);
    if (taskToArchive) {
      const archivedTask = {
        ...taskToArchive,
        completedAt: new Date().toISOString()
      };
      setTaskHistory(prev => [...prev, archivedTask]);
      setTasks(prev => prev.filter(task => task.id !== taskId));
    }
  };

  const restoreTask = (taskId: number) => {
    const taskToRestore = taskHistory.find(task => task.id === taskId);
    if (taskToRestore) {
      const restoredTask = {
        ...taskToRestore,
        completedAt: undefined,
        isExpanded: false
      };
      setTasks(prev => [...prev, restoredTask]);
      setTaskHistory(prev => prev.filter(task => task.id !== taskId));
    }
  };

  const updateUserProfile = (name: string, email: string) => {
    if (user) {
      const updatedUser = { ...user, name, email };
      setUser(updatedUser);
      localStorage.setItem('taskSplitter_user', JSON.stringify(updatedUser));
    }
  };

  // Подсчитываем общий прогресс
  const allSubtasks = tasks.flatMap(task => task.subtasks);
  const completedCount = allSubtasks.filter(s => s.completed).length;
  const totalCount = allSubtasks.length;

  return (
    <div className="App">
      <div className="chat-container">
        <div className="chat-header">
          <div className="header-content">
            <div className="header-left">
              <h1>🤖 Проще</h1>
              <p>сделай проще</p>
            </div>
            <div className="header-right">
              {!isAuthenticated ? (
                <>
                  <button 
                    onClick={() => setShowLogin(true)} 
                    className="header-button"
                  >
                    🔑 Войти
                  </button>
                  <button 
                    onClick={() => setShowRegister(true)} 
                    className="header-button"
                  >
                    ✨ Регистрация
                  </button>
                </>
              ) : (
                <>
                  <button 
                    onClick={() => setShowHistory(!showHistory)} 
                    className="header-button"
                  >
                    📚 История ({taskHistory.length})
                  </button>
                  <button 
                    onClick={() => setShowProfile(!showProfile)} 
                    className="header-button profile-button"
                  >
                    {user?.avatar} {user?.name}
                  </button>
                  <button 
                    onClick={handleLogout} 
                    className="header-button"
                  >
                    🚪 Выйти
                  </button>
                  {tasks.length > 0 && (
                    <button onClick={clearAllTasks} className="clear-button">
                      🗑️ Архивировать все
                    </button>
                  )}
                </>
              )}
            </div>
          </div>
        </div>

        {isAuthenticated ? (
          <div className="kanban-board">
            <div className="tasks-column">
            <h3>📋 Задачи</h3>
            <div className="tasks-list">
              {tasks.map(task => (
                <div key={task.id} className="task-card">
                  <div className="task-header">
                    <div 
                      className="task-header-content"
                      onClick={() => toggleTaskExpansion(task.id)}
                    >
                      <span className="task-title">{task.title}</span>
                      <span className="expand-icon">
                        {task.isExpanded ? '🔽' : '▶️'}
                      </span>
                    </div>
                    <button 
                      onClick={(e) => {
                        e.stopPropagation();
                        archiveTask(task.id);
                      }}
                      className="archive-button"
                      title="Архивировать задачу"
                    >
                      📦
                    </button>
                  </div>
                  
                  {task.isExpanded && (
                    <div className="subtasks-preview">
                      <div className="subtasks-count">
                        {task.subtasks.filter(s => s.completed).length} / {task.subtasks.length} выполнено
                      </div>
                      <div className="subtasks-list">
                        {task.subtasks.map(subtask => (
                          <div key={subtask.id} className={`subtask ${subtask.completed ? 'completed' : ''}`}>
                            <label className="subtask-item">
                              <input
                                type="checkbox"
                                checked={subtask.completed}
                                onChange={() => toggleSubtask(task.id, subtask.id)}
                              />
                              <span className="subtask-text">{subtask.title}</span>
                            </label>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>

          <div className="progress-column">
            <h3>📊 Прогресс</h3>
            {totalCount > 0 && (
              <div className="progress-section">
                <div className="progress-bar">
                  <div className="progress-header">
                    <div className="progress-text">
                      Выполнено: {completedCount} из {totalCount}
                    </div>
                    <div className="progress-track">
                      <div className="progress-fill" style={{ width: `${totalCount > 0 ? (completedCount / totalCount) * 100 : 0}%` }}></div>
                    </div>
                  </div>
                </div>
                
                <div className="tasks-summary">
                  {tasks.map(task => {
                    const taskCompleted = task.subtasks.filter(s => s.completed).length;
                    const taskTotal = task.subtasks.length;
                    return (
                      <div key={task.id} className="task-summary">
                        <span className="task-name">{task.title}</span>
                        <span className="task-progress">{taskCompleted}/{taskTotal}</span>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </div>
        </div>
        ) : (
          <div className="welcome-screen">
            <div className="welcome-content">
              <h2>🤖 Добро пожаловать в Проще!</h2>
              <p>Войдите или зарегистрируйтесь, чтобы начать разбивать задачи на простые шаги</p>
            </div>
          </div>
        )}

        {error && (
          <div className="error-message">
            ❌ {error}
          </div>
        )}

        {isAuthenticated && (
          <form onSubmit={handleSubmit} className="chat-input">
            <div className="input-group">
              <input
                type="text"
                value={task}
                onChange={(e) => setTask(e.target.value)}
                placeholder="Напишите задачу... например: помыть посуду"
                disabled={isLoading}
                className="task-input"
              />
              <button 
                type="submit" 
                disabled={isLoading || !task.trim()}
                className="submit-button"
              >
                {isLoading ? '⏳' : '🚀'}
              </button>
            </div>
          </form>
        )}

        {/* Модальное окно профиля */}
        {showProfile && (
          <div className="modal-overlay" onClick={() => setShowProfile(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>👤 Профиль пользователя</h3>
                <button onClick={() => setShowProfile(false)} className="close-button">✕</button>
              </div>
              <div className="modal-content">
                <div className="profile-info">
                  <div className="avatar-large">{user?.avatar}</div>
                  <div className="profile-details">
                    <p><strong>Имя:</strong> {user?.name}</p>
                    <p><strong>Email:</strong> {user?.email}</p>
                    <p><strong>Дата регистрации:</strong> {user?.createdAt ? new Date(user.createdAt).toLocaleDateString() : ''}</p>
                    <p><strong>Всего задач:</strong> {tasks.length + taskHistory.length}</p>
                    <p><strong>Активных задач:</strong> {tasks.length}</p>
                    <p><strong>В архиве:</strong> {taskHistory.length}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Модальное окно входа */}
        {showLogin && (
          <div className="modal-overlay" onClick={() => setShowLogin(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>🔑 Вход</h3>
                <button onClick={() => setShowLogin(false)} className="close-button">✕</button>
              </div>
              <div className="modal-content">
                <form onSubmit={handleLogin} className="register-form">
                  <div className="form-group">
                    <label htmlFor="loginUsername">Имя пользователя *</label>
                    <input
                      type="text"
                      id="loginUsername"
                      value={loginForm.username}
                      onChange={(e) => setLoginForm(prev => ({ ...prev, username: e.target.value }))}
                      placeholder="Введите имя пользователя"
                      required
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="loginPassword">Пароль *</label>
                    <input
                      type="password"
                      id="loginPassword"
                      value={loginForm.password}
                      onChange={(e) => setLoginForm(prev => ({ ...prev, password: e.target.value }))}
                      placeholder="Введите пароль"
                      required
                    />
                  </div>
                  
                  <button 
                    type="submit" 
                    disabled={isLoading}
                    className="register-button"
                  >
                    {isLoading ? '⏳ Вход...' : '🔑 Войти'}
                  </button>
                  
                  <div className="auth-switch">
                    <p>Нет аккаунта? <button type="button" onClick={() => { setShowLogin(false); setShowRegister(true); }} className="link-button">Зарегистрироваться</button></p>
                  </div>
                </form>
              </div>
            </div>
          </div>
        )}

        {/* Модальное окно регистрации */}
        {showRegister && (
          <div className="modal-overlay" onClick={() => setShowRegister(false)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>✨ Регистрация</h3>
                <button onClick={() => setShowRegister(false)} className="close-button">✕</button>
              </div>
              <div className="modal-content">
                <form onSubmit={handleRegister} className="register-form">
                  <div className="form-group">
                    <label htmlFor="username">Имя пользователя *</label>
                    <input
                      type="text"
                      id="username"
                      value={registerForm.username}
                      onChange={(e) => setRegisterForm(prev => ({ ...prev, username: e.target.value }))}
                      placeholder="Введите имя пользователя"
                      required
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="email">Email *</label>
                    <input
                      type="email"
                      id="email"
                      value={registerForm.email}
                      onChange={(e) => setRegisterForm(prev => ({ ...prev, email: e.target.value }))}
                      placeholder="Введите email"
                      required
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="password">Пароль *</label>
                    <input
                      type="password"
                      id="password"
                      value={registerForm.password}
                      onChange={(e) => setRegisterForm(prev => ({ ...prev, password: e.target.value }))}
                      placeholder="Введите пароль"
                      required
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="firstName">Имя</label>
                    <input
                      type="text"
                      id="firstName"
                      value={registerForm.firstName}
                      onChange={(e) => setRegisterForm(prev => ({ ...prev, firstName: e.target.value }))}
                      placeholder="Введите имя"
                    />
                  </div>
                  
                  <div className="form-group">
                    <label htmlFor="lastName">Фамилия</label>
                    <input
                      type="text"
                      id="lastName"
                      value={registerForm.lastName}
                      onChange={(e) => setRegisterForm(prev => ({ ...prev, lastName: e.target.value }))}
                      placeholder="Введите фамилию"
                    />
                  </div>
                  
                  <button 
                    type="submit" 
                    disabled={isLoading}
                    className="register-button"
                  >
                    {isLoading ? '⏳ Регистрация...' : '✨ Зарегистрироваться'}
                  </button>
                  
                  <div className="auth-switch">
                    <p>Уже есть аккаунт? <button type="button" onClick={() => { setShowRegister(false); setShowLogin(true); }} className="link-button">Войти</button></p>
                  </div>
                </form>
              </div>
            </div>
          </div>
        )}

        {/* Модальное окно истории */}
        {showHistory && (
          <div className="modal-overlay" onClick={() => setShowHistory(false)}>
            <div className="modal history-modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>📚 История задач</h3>
                <button onClick={() => setShowHistory(false)} className="close-button">✕</button>
              </div>
              <div className="modal-content">
                {taskHistory.length === 0 ? (
                  <div className="empty-history">
                    <p>История пуста</p>
                    <p>Завершенные задачи будут появляться здесь</p>
                  </div>
                ) : (
                  <div className="history-list">
                    {taskHistory.map(task => (
                      <div key={task.id} className="history-item">
                        <div className="history-task-info">
                          <h4>{task.title}</h4>
                          <p>Создано: {new Date(task.createdAt).toLocaleDateString()}</p>
                          <p>Завершено: {task.completedAt ? new Date(task.completedAt).toLocaleDateString() : 'Не завершено'}</p>
                          <p>Подзадач: {task.subtasks.length}</p>
                        </div>
                        <button 
                          onClick={() => restoreTask(task.id)}
                          className="restore-button"
                          title="Восстановить задачу"
                        >
                          🔄 Восстановить
                        </button>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;