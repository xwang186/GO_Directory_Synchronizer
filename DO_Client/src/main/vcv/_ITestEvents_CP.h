#pragma once

template<class T>
class CProxy_ITestEvents :
	public ATL::IConnectionPointImpl<T, &__uuidof(_ITestEvents)>
{
public:
	HRESULT Fire_Notify(BSTR str)
	{
    std::cout << "\n  -- component calling Fire_Notify ---";
		HRESULT hr = S_OK;
		T * pThis = static_cast<T *>(this);
		int cConnections = m_vec.GetSize();

		for (int iConnection = 0; iConnection < cConnections; iConnection++)
		{
			pThis->Lock();
			CComPtr<IUnknown> punkConnection = m_vec.GetAt(iConnection);
			pThis->Unlock();

			IDispatch * pConnection = static_cast<IDispatch *>(punkConnection.p);

			if (pConnection)
			{
        std::cout << "\n  -- Component Firing Notify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

tify ---";
		HRESULT hr = S_OK;
		T * pThis = static_cast<T *>(this);
		int cConnections = m_vec.GetSize();

		for (int iConnection = 0; iConnection < cConnections; iConnection++)
		{
			pThis->Lock();
			CComPtr<IUnknown> punkConnection = m_vec.GetAt(iConnection);
			pThis->Unlock();

			IDispatch * pConnection = static_cast<IDispatch *>(punkConnection.p);

			if (pConnection)
			{
        std::cout << "\n  -- Component Firing Notify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};


		for (int iConnection = 0; iConnection < cConnections; iConnection++)
		{
			pThis->Lock();
			CComPtr<IUnknown> punkConnection = m_vec.GetAt(iConnection);
			pThis->Unlock();

			IDispatch * pConnection = static_cast<IDispatch *>(punkConnection.p);

			if (pConnection)
			{
        std::cout << "\n  -- Component Firing Notify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

known> punkConnection = m_vec.GetAt(iConnection);
			pThis->Unlock();

			IDispatch * pConnection = static_cast<IDispatch *>(punkConnection.p);

			if (pConnection)
			{
        std::cout << "\n  -- Component Firing Notify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

st<IDispatch *>(punkConnection.p);

			if (pConnection)
			{
        std::cout << "\n  -- Component Firing Notify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

otify Event --";
				CComVariant avarParams[1];
				avarParams[0] = str;
				avarParams[0].vt = VT_BSTR;
				CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

	CComVariant varResult;
				DISPPARAMS params = { avarParams, NULL, 1, 0 };

				hr = pConnection->Invoke(1, IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

IID_NULL, LOCALE_USER_DEFAULT, DISPATCH_METHOD, &params, &varResult, NULL, NULL);
        if (!SUCCEEDED(hr))
        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

        {
          std::cout << "\n  -- callback failed --";
        }
      }
		}
		return hr;
	}
};

