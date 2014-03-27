<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Auth extends CI_Controller {

	function index(){
		redirectIfNoSession($this);
		$this->load->view('auth/login');
	}

	public function login(){
		//echo '<pre>'; print_r(apache_get_modules());echo '</pre>';exit;
		redirectIfSession($this);
		$this->load->view('auth/login');
	}

	public function register(){
		redirectIfSession($this);
		$this->load->view('auth/register');
	}

	public function doRegister(){
		redirectIfSession($this);

		$this->load->library('form_validation');

		$this->form_validation->set_rules('password', 'Password', 'required');
		$this->form_validation->set_rules('email', 'Email', 'required|valid_email');

		if ($this->form_validation->run() == FALSE){
			echo json_encode(array(
				'status' => 'error'
				,'msg' => $this->form_validation->error_string()
			));
			exit;
		}

		$post = $this->input->post();

		//lets load user model
		//u is alias of User_db
		$this->load->model('User_db','u');

		//lets create new user by inserting into db
		$user_id = $this->u->insertNewUser($post['email'],$post['password']);

		if($user_id <= 0){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> $user_id == 0 ? 'Somer error occured. Please try after some time.' : 'An user is already existing with this account. Please user correct password to login.'
			));
			exit;
		}

		$random_string = substr(sha1(rand()), 0, 15);

		//lets set session
		$this->session->set_userdata(array(
			'user_id' => $user_id
			,'email' => $post['email']
			,'token' => $random_string
		));

		//LETS REGISTER IN GO SERVER
		$status = post_to_go('auth','',array('USER_ID' => $user_id, 'TOKEN' => $random_string),array('USER_ID : '.$user_id, 'TOKEN : '.$random_string));

		echo json_encode(array(
			'status' => 'success'
			,'message' => 'You are successfully registered.'
			,'user_id' => $user_id
			,'token' => $random_string
		));
		exit;
	}

	function doLogin(){
		redirectIfSession($this);

		$this->load->library('form_validation');

		$this->form_validation->set_rules('password', 'Password', 'required');
		$this->form_validation->set_rules('email', 'Email', 'required|valid_email');

		if ($this->form_validation->run() == FALSE){
			echo json_encode(array(
				'status' => 'error'
				,'msg' => $this->form_validation->error_string()
			));
			exit;
		}

		$post = $this->input->post();

		//lets load user model
		//u is alias of User_db
		$this->load->model('User_db','u');

		//lets get user details from db
		$data = $this->u->getUserDetails($post['email']);

		if(!is_object($data)){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> 'This email is not registered with us. Please register'
			));
			exit;
		}

		if($data->password != $post['password']){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> 'Please provide correct password'
			));
			exit;
		}

		$random_string = substr(sha1(rand()), 0, 15);

		//lets set session
		$this->session->set_userdata(array(
			'user_id' => $data->user_id
			,'email' => $post['email']
			,'token' => $random_string
		));

		//LETS REGISTER IN GO SERVER
		$status = post_to_go('auth','',array('USER_ID' => $data->user_id, 'TOKEN' => $random_string),array('USER_ID: '.$data->user_id, 'TOKEN: '.$random_string));

		echo json_encode(array(
			'status' => 'success'
			,'message' => 'You are successfully logged-in.'
			,'user_id' => $user_id
			,'token' => $random_string
			,'result' => $status
		));
		exit;
	}

	function logout(){
		$this->session->destroy();
		redirect(base_url());
	}

	function test(){
		//post_to_go('add','user',array(),array('USER_ID' => $user_id, 'TOKEN' => $random_string));
	}
}